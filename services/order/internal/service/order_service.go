package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/order/internal/domain"
	"github.com/gondolia/gondolia/services/order/internal/repository"
)

// CartResponse represents the cart data from Cart Service
type CartResponse struct {
	ID     uuid.UUID      `json:"id"`
	Items  []CartItem     `json:"items"`
	Total  float64        `json:"total"`
}

type CartItem struct {
	ID            uuid.UUID      `json:"id"`
	ProductID     uuid.UUID      `json:"product_id"`
	VariantID     *uuid.UUID     `json:"variant_id"`
	ProductType   string         `json:"product_type"`
	ProductName   string         `json:"product_name"`
	SKU           string         `json:"sku"`
	Quantity      int            `json:"quantity"`
	UnitPrice     float64        `json:"unit_price"`
	Currency      string         `json:"currency"`
	Configuration map[string]any `json:"configuration"`
}

type OrderService struct {
	orderRepo      repository.OrderRepository
	tenantRepo     repository.TenantRepository
	cartServiceURL string
	httpClient     *http.Client
}

func NewOrderService(orderRepo repository.OrderRepository, tenantRepo repository.TenantRepository, cartServiceURL string) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		tenantRepo:     tenantRepo,
		cartServiceURL: cartServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Checkout creates an order from a cart
func (s *OrderService) Checkout(ctx context.Context, tenantID, userID uuid.UUID, sessionID string, req *domain.CheckoutRequest) (*domain.Order, error) {
	// 1. Get cart from Cart Service
	cart, err := s.getCart(ctx, tenantID, userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, domain.ErrCartEmpty
	}

	// 2. Validate cart (call Cart Service validate endpoint)
	if err := s.validateCart(ctx, tenantID, userID, sessionID); err != nil {
		return nil, fmt.Errorf("cart validation failed: %w", err)
	}

	// 3. Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}

	// 4. Create order
	order := domain.NewOrder(tenantID, userID, orderNumber)
	order.ShippingAddress = req.ShippingAddress
	order.BillingAddress = req.BillingAddress
	order.Notes = req.Notes

	// Calculate totals from cart items
	var subtotal float64
	orderItems := make([]domain.OrderItem, 0, len(cart.Items))

	for _, cartItem := range cart.Items {
		itemTotal := cartItem.UnitPrice * float64(cartItem.Quantity)
		subtotal += itemTotal

		orderItem := domain.OrderItem{
			ID:            uuid.New(),
			OrderID:       order.ID,
			ProductID:     cartItem.ProductID,
			VariantID:     cartItem.VariantID,
			ProductType:   cartItem.ProductType,
			ProductName:   cartItem.ProductName,
			SKU:           cartItem.SKU,
			Quantity:      cartItem.Quantity,
			UnitPrice:     cartItem.UnitPrice,
			TotalPrice:    itemTotal,
			Currency:      cartItem.Currency,
			Configuration: cartItem.Configuration,
			CreatedAt:     time.Now(),
		}
		orderItems = append(orderItems, orderItem)
	}

	// Set order totals (simplified - no tax calculation for now)
	order.Subtotal = subtotal
	order.TaxAmount = 0 // TODO: Implement tax calculation
	order.Total = subtotal
	order.Currency = "EUR" // Default currency
	if len(cart.Items) > 0 {
		order.Currency = cart.Items[0].Currency
	}

	// 5. Save order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 6. Save order items
	if err := s.orderRepo.CreateItems(ctx, orderItems); err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// 7. Add status history log
	statusLog := &domain.OrderStatusLog{
		ID:        uuid.New(),
		OrderID:   order.ID,
		FromStatus: nil,
		ToStatus:  order.Status,
		Note:      "Order created from checkout",
		CreatedAt: time.Now(),
	}
	if err := s.orderRepo.AddStatusLog(ctx, statusLog); err != nil {
		return nil, fmt.Errorf("failed to add status log: %w", err)
	}

	// 8. Mark cart as completed
	if err := s.markCartCompleted(ctx, tenantID, userID, sessionID); err != nil {
		// Log error but don't fail the order
		fmt.Printf("Warning: failed to mark cart as completed: %v\n", err)
	}

	// Load full order with items and history
	order.Items = orderItems
	order.StatusHistory = []domain.OrderStatusLog{*statusLog}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return s.orderRepo.GetByID(ctx, orderID)
}

// ListOrders lists orders for a user
func (s *OrderService) ListOrders(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int, error) {
	return s.orderRepo.List(ctx, filter)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID uuid.UUID) (*domain.Order, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Check if order belongs to user
	if order.UserID != userID {
		return nil, domain.ErrForbidden
	}

	// Check if order can be cancelled
	if !order.CanBeCancelled() {
		return nil, domain.ErrOrderCannotBeCancelled
	}

	// Update status
	oldStatus := order.Status
	order.Status = domain.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Add status history log
	statusLog := &domain.OrderStatusLog{
		ID:         uuid.New(),
		OrderID:    order.ID,
		FromStatus: &oldStatus,
		ToStatus:   domain.OrderStatusCancelled,
		ChangedBy:  &userID,
		Note:       "Order cancelled by user",
		CreatedAt:  time.Now(),
	}
	if err := s.orderRepo.AddStatusLog(ctx, statusLog); err != nil {
		return nil, err
	}

	return order, nil
}

// getCart retrieves cart from Cart Service
func (s *OrderService) getCart(ctx context.Context, tenantID, userID uuid.UUID, sessionID string) (*CartResponse, error) {
	// Get tenant code from database
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/cart", s.cartServiceURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Use tenant CODE instead of UUID
	req.Header.Set("X-Tenant-ID", tenant.Code)
	req.Header.Set("X-User-ID", userID.String())
	if sessionID != "" {
		req.Header.Set("X-Session-ID", sessionID)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cart service returned status %d: %s", resp.StatusCode, string(body))
	}

	var cart CartResponse
	if err := json.NewDecoder(resp.Body).Decode(&cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

// validateCart validates cart items with Cart Service
func (s *OrderService) validateCart(ctx context.Context, tenantID, userID uuid.UUID, sessionID string) error {
	// Get tenant code from database
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/cart/validate", s.cartServiceURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	// Use tenant CODE instead of UUID
	req.Header.Set("X-Tenant-ID", tenant.Code)
	req.Header.Set("X-User-ID", userID.String())
	if sessionID != "" {
		req.Header.Set("X-Session-ID", sessionID)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cart validation failed: %s", string(body))
	}

	return nil
}

// markCartCompleted marks cart as completed
func (s *OrderService) markCartCompleted(ctx context.Context, tenantID, userID uuid.UUID, sessionID string) error {
	// Get tenant code from database
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/cart/complete", s.cartServiceURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// Use tenant CODE instead of UUID
	req.Header.Set("X-Tenant-ID", tenant.Code)
	req.Header.Set("X-User-ID", userID.String())
	if sessionID != "" {
		req.Header.Set("X-Session-ID", sessionID)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to mark cart as completed: %s", string(body))
	}

	return nil
}
