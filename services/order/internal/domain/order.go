package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID              uuid.UUID         `json:"id"`
	TenantID        uuid.UUID         `json:"tenant_id"`
	UserID          uuid.UUID         `json:"user_id"`
	OrderNumber     string            `json:"order_number"`
	Status          OrderStatus       `json:"status"`
	Subtotal        float64           `json:"subtotal"`
	TaxAmount       float64           `json:"tax_amount"`
	Total           float64           `json:"total"`
	Currency        string            `json:"currency"`
	ShippingAddress map[string]any    `json:"shipping_address"`
	BillingAddress  map[string]any    `json:"billing_address"`
	Notes           string            `json:"notes,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ItemCount       int               `json:"item_count,omitempty"`       // Number of items (for list views)
	Items           []OrderItem       `json:"items,omitempty"`
	StatusHistory   []OrderStatusLog  `json:"status_history,omitempty"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID            uuid.UUID      `json:"id"`
	OrderID       uuid.UUID      `json:"order_id"`
	ProductID     uuid.UUID      `json:"product_id"`
	VariantID     *uuid.UUID     `json:"variant_id,omitempty"`
	ProductType   string         `json:"product_type"`
	ProductName   string         `json:"product_name"`
	SKU           string         `json:"sku"`
	Quantity      int            `json:"quantity"`
	UnitPrice     float64        `json:"unit_price"`
	TotalPrice    float64        `json:"total_price"`
	Currency      string         `json:"currency"`
	Configuration map[string]any `json:"configuration,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}

// OrderStatusLog represents a status change in the order history
type OrderStatusLog struct {
	ID         uuid.UUID   `json:"id"`
	OrderID    uuid.UUID   `json:"order_id"`
	FromStatus *OrderStatus `json:"from_status,omitempty"`
	ToStatus   OrderStatus `json:"to_status"`
	ChangedBy  *uuid.UUID  `json:"changed_by,omitempty"`
	Note       string      `json:"note,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

// CheckoutRequest represents a request to create an order from a cart
type CheckoutRequest struct {
	ShippingAddress map[string]any `json:"shipping_address" binding:"required"`
	BillingAddress  map[string]any `json:"billing_address" binding:"required"`
	Notes           string         `json:"notes,omitempty"`
}

// OrderFilter represents filter options for listing orders
type OrderFilter struct {
	TenantID uuid.UUID
	UserID   *uuid.UUID
	Status   *OrderStatus
	Limit    int
	Offset   int
}

// NewOrder creates a new order with defaults
func NewOrder(tenantID, userID uuid.UUID, orderNumber string) *Order {
	now := time.Now()
	return &Order{
		ID:          uuid.New(),
		TenantID:    tenantID,
		UserID:      userID,
		OrderNumber: orderNumber,
		Status:      OrderStatusConfirmed,
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       []OrderItem{},
	}
}

// CanBeCancelled checks if the order can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

// IsValidStatusTransition checks if a status transition is valid
func IsValidStatusTransition(from, to OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
		OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:    {OrderStatusDelivered},
		OrderStatusDelivered:  {},
		OrderStatusCancelled:  {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == to {
			return true
		}
	}
	return false
}
