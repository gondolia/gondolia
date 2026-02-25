package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/cart/internal/domain"
	"github.com/gondolia/gondolia/services/cart/internal/repository"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo          repository.CartRepository
	tenantRepo        repository.TenantRepository
	catalogServiceURL string
}

// NewCartService creates a new cart service
func NewCartService(cartRepo repository.CartRepository, tenantRepo repository.TenantRepository, catalogServiceURL string) *CartService {
	return &CartService{
		cartRepo:          cartRepo,
		tenantRepo:        tenantRepo,
		catalogServiceURL: catalogServiceURL,
	}
}

// GetOrCreateCart gets the active cart or creates a new one
func (s *CartService) GetOrCreateCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) (*domain.Cart, error) {
	// Try to get existing active cart
	cart, err := s.cartRepo.GetActiveCart(ctx, tenantID, userID, sessionID)
	if err == nil {
		return cart, nil
	}

	// If not found, create new cart
	if err == domain.ErrCartNotFound {
		cart = domain.NewCart(tenantID, userID, sessionID)
		if err := s.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
		return cart, nil
	}

	return nil, err
}

// GetCart retrieves cart by ID
func (s *CartService) GetCart(ctx context.Context, id uuid.UUID) (*domain.Cart, error) {
	return s.cartRepo.GetByID(ctx, id)
}

// AddItem adds an item to the cart, or increases quantity if an identical item exists
func (s *CartService) AddItem(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string, req domain.AddItemRequest) (*domain.Cart, error) {
	// Get or create cart
	cart, err := s.GetOrCreateCart(ctx, tenantID, userID, sessionID)
	if err != nil {
		return nil, err
	}

	if !cart.IsActive() {
		return nil, domain.ErrCartNotActive
	}

	// Check if an identical item already exists in the cart
	// (same product_id, variant_id, and configuration)
	configHash := computeConfigHash(req.Configuration)
	existingItem, err := s.cartRepo.FindMatchingItem(ctx, cart.ID, req.ProductID, req.VariantID, configHash)

	now := time.Now()

	if err == nil && existingItem != nil {
		// Item exists - increase quantity
		newQuantity := existingItem.Quantity + req.Quantity
		existingItem.Quantity = newQuantity
		existingItem.TotalPrice = existingItem.UnitPrice * float64(newQuantity)
		existingItem.UpdatedAt = now

		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		// No existing item - create new one
		// Fetch product metadata and price from catalog service
		productInfo, err := s.fetchProductInfo(ctx, tenantID, req.ProductID, req.VariantID, req.Quantity, req.Configuration)
		if err != nil {
			return nil, err
		}

		// Create cart item
		item := &domain.CartItem{
			ID:            uuid.New(),
			CartID:        cart.ID,
			ProductID:     req.ProductID,
			VariantID:     req.VariantID,
			ProductType:   productInfo.ProductType,
			ProductName:   productInfo.Name,
			SKU:           productInfo.SKU,
			ImageURL:      productInfo.ImageURL,
			Quantity:      req.Quantity,
			UnitPrice:     productInfo.UnitPrice,
			TotalPrice:    productInfo.UnitPrice * float64(req.Quantity),
			Currency:      productInfo.Currency,
			Configuration: req.Configuration,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Add item to database
		if err := s.cartRepo.AddItem(ctx, item); err != nil {
			return nil, err
		}
	}

	// Update cart timestamp
	cart.UpdatedAt = now
	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	// Reload cart with items
	return s.cartRepo.GetByID(ctx, cart.ID)
}

// computeConfigHash computes an MD5 hash of the configuration for comparison
// Returns empty string for nil configuration
func computeConfigHash(config *domain.Configuration) string {
	if config == nil {
		return ""
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", md5.Sum(configJSON))
}

// UpdateItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) (*domain.Cart, error) {
	if quantity < 1 {
		return nil, domain.ErrInvalidQuantity
	}

	// Get item
	item, err := s.cartRepo.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Update quantity and recalculate total price
	item.Quantity = quantity
	item.TotalPrice = item.UnitPrice * float64(quantity)
	item.UpdatedAt = time.Now()

	if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	// Get cart and update timestamp
	cart, err := s.cartRepo.GetByID(ctx, item.CartID)
	if err != nil {
		return nil, err
	}

	cart.UpdatedAt = time.Now()
	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	// Reload cart with items
	return s.cartRepo.GetByID(ctx, cart.ID)
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(ctx context.Context, itemID uuid.UUID) (*domain.Cart, error) {
	// Get item first to get cart ID
	item, err := s.cartRepo.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(ctx, itemID); err != nil {
		return nil, err
	}

	// Get cart and update timestamp
	cart, err := s.cartRepo.GetByID(ctx, item.CartID)
	if err != nil {
		return nil, err
	}

	cart.UpdatedAt = time.Now()
	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	// Reload cart with items
	return s.cartRepo.GetByID(ctx, cart.ID)
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) error {
	// Get active cart
	cart, err := s.cartRepo.GetActiveCart(ctx, tenantID, userID, sessionID)
	if err != nil {
		return err
	}

	// Clear all items
	if err := s.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return err
	}

	// Update cart timestamp
	cart.UpdatedAt = time.Now()
	return s.cartRepo.Update(ctx, cart)
}

// ValidateCart validates all items in the cart against catalog
func (s *CartService) ValidateCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) (*domain.Cart, error) {
	// Get active cart
	cart, err := s.cartRepo.GetActiveCart(ctx, tenantID, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// Validate each item and update prices
	for i, item := range cart.Items {
		productInfo, err := s.fetchProductInfo(ctx, tenantID, item.ProductID, item.VariantID, item.Quantity, item.Configuration)
		if err != nil {
			// Item no longer available
			continue
		}

		// Update if price or product info changed
		needsUpdate := false
		if item.UnitPrice != productInfo.UnitPrice || item.Currency != productInfo.Currency {
			cart.Items[i].UnitPrice = productInfo.UnitPrice
			cart.Items[i].TotalPrice = productInfo.UnitPrice * float64(item.Quantity)
			cart.Items[i].Currency = productInfo.Currency
			needsUpdate = true
		}
		if item.ProductName != productInfo.Name || item.SKU != productInfo.SKU || item.ImageURL != productInfo.ImageURL {
			cart.Items[i].ProductName = productInfo.Name
			cart.Items[i].SKU = productInfo.SKU
			cart.Items[i].ImageURL = productInfo.ImageURL
			needsUpdate = true
		}

		if needsUpdate {
			cart.Items[i].UpdatedAt = time.Now()
			if err := s.cartRepo.UpdateItem(ctx, &cart.Items[i]); err != nil {
				return nil, err
			}
		}
	}

	return cart, nil
}

// MergeCarts merges a guest cart into a user cart (called on login)
func (s *CartService) MergeCarts(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, sessionID string) (*domain.Cart, error) {
	// Get guest cart
	guestCart, err := s.cartRepo.GetActiveCart(ctx, tenantID, nil, &sessionID)
	if err != nil {
		if err == domain.ErrCartNotFound {
			// No guest cart, just get or create user cart
			return s.GetOrCreateCart(ctx, tenantID, &userID, nil)
		}
		return nil, err
	}

	// Get or create user cart
	userCart, err := s.GetOrCreateCart(ctx, tenantID, &userID, nil)
	if err != nil {
		return nil, err
	}

	// Merge carts
	if err := s.cartRepo.MergeCarts(ctx, guestCart.ID, userCart.ID); err != nil {
		return nil, err
	}

	// Reload user cart
	return s.cartRepo.GetByID(ctx, userCart.ID)
}

// CompleteCart marks a cart as completed (after order creation)
func (s *CartService) CompleteCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) error {
	// Get active cart
	cart, err := s.cartRepo.GetActiveCart(ctx, tenantID, userID, sessionID)
	if err != nil {
		return err
	}

	// Mark as completed
	cart.Status = domain.CartStatusCompleted
	cart.UpdatedAt = time.Now()
	return s.cartRepo.Update(ctx, cart)
}

// ProductInfo holds product metadata and pricing information
type ProductInfo struct {
	ProductType domain.ProductType
	Name        string
	SKU         string
	ImageURL    string
	UnitPrice   float64
	Currency    string
}

// fetchProductInfo fetches product metadata and price from catalog service
func (s *CartService) fetchProductInfo(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, quantity int, configuration *domain.Configuration) (*ProductInfo, error) {
	// Fetch price using existing function
	price, currency, productType, err := s.fetchPriceFromCatalog(ctx, tenantID, productID, variantID, quantity, configuration)
	if err != nil {
		return nil, err
	}

	// Fetch product metadata
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	lookupID := productID
	if variantID != nil {
		lookupID = *variantID
	}

	client := &http.Client{Timeout: 10 * time.Second}
	productURL := fmt.Sprintf("%s/api/v1/products/%s", s.catalogServiceURL, lookupID)
	productReq, err := http.NewRequestWithContext(ctx, "GET", productURL, nil)
	if err != nil {
		return nil, err
	}
	productReq.Header.Set("X-Tenant-ID", tenant.Code)

	productResp, err := client.Do(productReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product from catalog: %w", err)
	}
	defer productResp.Body.Close()

	if productResp.StatusCode != http.StatusOK {
		return nil, domain.ErrProductNotFound
	}

	var productData struct {
		Data struct {
			Name     json.RawMessage `json:"name"`
			SKU      string          `json:"sku"`
			ImageURL string          `json:"image_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(productResp.Body).Decode(&productData); err != nil {
		return nil, fmt.Errorf("failed to parse product response: %w", err)
	}

	// Parse name: could be a string or an i18n object {"de":"...", "en":"..."}
	var productName string
	if err := json.Unmarshal(productData.Data.Name, &productName); err != nil {
		// Try as i18n map
		var nameMap map[string]string
		if err2 := json.Unmarshal(productData.Data.Name, &nameMap); err2 == nil {
			if de, ok := nameMap["de"]; ok {
				productName = de
			} else if en, ok := nameMap["en"]; ok {
				productName = en
			} else {
				for _, v := range nameMap {
					productName = v
					break
				}
			}
		}
	}

	return &ProductInfo{
		ProductType: productType,
		Name:        productName,
		SKU:         productData.Data.SKU,
		ImageURL:    productData.Data.ImageURL,
		UnitPrice:   price,
		Currency:    currency,
	}, nil
}

// fetchPriceFromCatalog fetches product price from catalog service via HTTP
// Handles different pricing strategies based on product type:
// - simple/variant: fetches tiered prices and selects appropriate tier based on quantity
// - bundle: calculates price based on configuration
// - parametric: calculates price based on parameters and selections
func (s *CartService) fetchPriceFromCatalog(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, quantity int, configuration *domain.Configuration) (float64, string, domain.ProductType, error) {
	// Get tenant code from database
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to get tenant: %w", err)
	}

	// Determine which ID to use for lookups
	lookupID := productID
	if variantID != nil {
		lookupID = *variantID
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Get product info to determine product_type
	productURL := fmt.Sprintf("%s/api/v1/products/%s", s.catalogServiceURL, lookupID)
	productReq, err := http.NewRequestWithContext(ctx, "GET", productURL, nil)
	if err != nil {
		return 0, "", "", err
	}
	productReq.Header.Set("X-Tenant-ID", tenant.Code)

	productResp, err := client.Do(productReq)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to fetch product from catalog: %w", err)
	}
	defer productResp.Body.Close()

	if productResp.StatusCode != http.StatusOK {
		return 0, "", "", domain.ErrProductNotFound
	}

	var productData struct {
		Data struct {
			ProductType string `json:"product_type"`
		} `json:"data"`
	}

	if err := json.NewDecoder(productResp.Body).Decode(&productData); err != nil {
		return 0, "", "", fmt.Errorf("failed to parse product response: %w", err)
	}

	productType := domain.ProductType(productData.Data.ProductType)

	// Step 2: Fetch price based on product_type
	switch productType {
	case domain.ProductTypeSimple, domain.ProductTypeVariant:
		// Fetch tiered prices
		pricesURL := fmt.Sprintf("%s/api/v1/products/%s/prices", s.catalogServiceURL, lookupID)
		pricesReq, err := http.NewRequestWithContext(ctx, "GET", pricesURL, nil)
		if err != nil {
			return 0, "", "", err
		}
		pricesReq.Header.Set("X-Tenant-ID", tenant.Code)

		pricesResp, err := client.Do(pricesReq)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to fetch prices from catalog: %w", err)
		}
		defer pricesResp.Body.Close()

		if pricesResp.StatusCode != http.StatusOK {
			return 0, "", "", domain.ErrPriceNotAvailable
		}

		var pricesData struct {
			Data []struct {
				MinQuantity int     `json:"min_quantity"`
				Price       float64 `json:"price"`
				Currency    string  `json:"currency"`
			} `json:"data"`
		}

		if err := json.NewDecoder(pricesResp.Body).Decode(&pricesData); err != nil {
			return 0, "", "", fmt.Errorf("failed to parse prices response: %w", err)
		}

		if len(pricesData.Data) == 0 {
			return 0, "", "", domain.ErrPriceNotAvailable
		}

		// Find the appropriate tier price based on quantity
		// Prices are typically sorted by min_quantity, but we'll be safe
		var selectedPrice *struct {
			MinQuantity int     `json:"min_quantity"`
			Price       float64 `json:"price"`
			Currency    string  `json:"currency"`
		}

		for i := range pricesData.Data {
			tier := &pricesData.Data[i]
			if tier.MinQuantity <= quantity {
				if selectedPrice == nil || tier.MinQuantity > selectedPrice.MinQuantity {
					selectedPrice = tier
				}
			}
		}

		// If no tier matches, use the lowest min_quantity tier
		if selectedPrice == nil {
			selectedPrice = &pricesData.Data[0]
			for i := range pricesData.Data {
				if pricesData.Data[i].MinQuantity < selectedPrice.MinQuantity {
					selectedPrice = &pricesData.Data[i]
				}
			}
		}

		return selectedPrice.Price, selectedPrice.Currency, productType, nil

	case domain.ProductTypeBundle:
		// Calculate bundle price with configuration
		bundleURL := fmt.Sprintf("%s/api/v1/bundles/%s/calculate-price", s.catalogServiceURL, lookupID)

		// Transform configuration from frontend format to catalog API format:
		// Frontend sends: { bundleComponents: [{ componentId, quantity, ... }] }
		// Catalog expects: { components: [{ component_id, quantity, ... }] }
		catalogConfig := transformBundleConfigForCatalog(configuration)

		configJSON, err := json.Marshal(catalogConfig)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to marshal bundle configuration: %w", err)
		}

		bundleReq, err := http.NewRequestWithContext(ctx, "POST", bundleURL, bytes.NewReader(configJSON))
		if err != nil {
			return 0, "", "", err
		}
		bundleReq.Header.Set("X-Tenant-ID", tenant.Code)
		bundleReq.Header.Set("Content-Type", "application/json")

		bundleResp, err := client.Do(bundleReq)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to calculate bundle price: %w", err)
		}
		defer bundleResp.Body.Close()

		if bundleResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(bundleResp.Body)
			return 0, "", "", fmt.Errorf("bundle price calculation failed with status %d: %s", bundleResp.StatusCode, string(body))
		}

		var bundleData struct {
			Total    float64 `json:"total"`
			Currency string  `json:"currency"`
		}

		if err := json.NewDecoder(bundleResp.Body).Decode(&bundleData); err != nil {
			return 0, "", "", fmt.Errorf("failed to parse bundle price response: %w", err)
		}

		// Fallback to CHF if currency is empty
		currency := bundleData.Currency
		if currency == "" {
			currency = "CHF"
		}

		return bundleData.Total, currency, productType, nil

	case domain.ProductTypeParametric:
		// Calculate parametric price with parameters and selections
		parametricURL := fmt.Sprintf("%s/api/v1/products/%s/calculate-price", s.catalogServiceURL, lookupID)

		// Transform configuration to catalog API format:
		// Catalog expects: { selections: {...}, parameters: {...}, quantity: N }
		catalogConfig := transformParametricConfigForCatalog(configuration, quantity)

		configJSON, err := json.Marshal(catalogConfig)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to marshal parametric configuration: %w", err)
		}

		parametricReq, err := http.NewRequestWithContext(ctx, "POST", parametricURL, bytes.NewReader(configJSON))
		if err != nil {
			return 0, "", "", err
		}
		parametricReq.Header.Set("X-Tenant-ID", tenant.Code)
		parametricReq.Header.Set("Content-Type", "application/json")

		parametricResp, err := client.Do(parametricReq)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to calculate parametric price: %w", err)
		}
		defer parametricResp.Body.Close()

		if parametricResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(parametricResp.Body)
			return 0, "", "", fmt.Errorf("parametric price calculation failed with status %d: %s", parametricResp.StatusCode, string(body))
		}

		var parametricData struct {
			UnitPrice float64 `json:"unit_price"`
			Currency  string  `json:"currency"`
		}

		if err := json.NewDecoder(parametricResp.Body).Decode(&parametricData); err != nil {
			return 0, "", "", fmt.Errorf("failed to parse parametric price response: %w", err)
		}

		return parametricData.UnitPrice, parametricData.Currency, productType, nil

	default:
		return 0, "", "", fmt.Errorf("unknown product type: %s", productType)
	}
}

// CatalogBundleConfig represents the format expected by the catalog service
type CatalogBundleConfig struct {
	Components []CatalogBundleComponent `json:"components"`
}

// CatalogBundleComponent represents a component in catalog API format
type CatalogBundleComponent struct {
	ComponentID string                 `json:"component_id"`
	Quantity    int                    `json:"quantity"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Selections  map[string]string      `json:"selections,omitempty"`
}

// transformBundleConfigForCatalog transforms frontend bundle configuration to catalog API format
// Frontend sends: { bundleComponents: [{ componentId, quantity, parameters, selections }] }
// Catalog expects: { components: [{ component_id, quantity, parameters, selections }] }
func transformBundleConfigForCatalog(config *domain.Configuration) *CatalogBundleConfig {
	if config == nil {
		return &CatalogBundleConfig{Components: []CatalogBundleComponent{}}
	}

	catalogConfig := &CatalogBundleConfig{
		Components: make([]CatalogBundleComponent, 0, len(config.BundleComponents)),
	}

	for _, comp := range config.BundleComponents {
		catalogComp := CatalogBundleComponent{
			ComponentID: comp.ComponentID.String(),
			Quantity:    comp.Quantity,
			Parameters:  comp.Parameters,
			Selections:  comp.Selections,
		}

		catalogConfig.Components = append(catalogConfig.Components, catalogComp)
	}

	return catalogConfig
}

// CatalogParametricConfig represents the format expected by the catalog service for parametric products
type CatalogParametricConfig struct {
	Selections map[string]string      `json:"selections"`
	Parameters map[string]interface{} `json:"parameters"`
	Quantity   int                    `json:"quantity"`
}

// transformParametricConfigForCatalog transforms frontend parametric configuration to catalog API format
// Frontend sends: { parameters: {...}, selections: {...} }
// Catalog expects: { selections: {...}, parameters: {...}, quantity: N }
func transformParametricConfigForCatalog(config *domain.Configuration, quantity int) *CatalogParametricConfig {
	catalogConfig := &CatalogParametricConfig{
		Selections: make(map[string]string),
		Parameters: make(map[string]interface{}),
		Quantity:   quantity,
	}

	if config == nil {
		return catalogConfig
	}

	// Copy parametric params from configuration
	if config.ParametricParams != nil {
		// The frontend sends parameters under "parameters" and selections under "selections"
		// within ParametricParams or as separate fields
		if params, ok := config.ParametricParams["parameters"]; ok {
			if paramsMap, ok := params.(map[string]interface{}); ok {
				catalogConfig.Parameters = paramsMap
			}
		} else {
			// Parameters might be directly in ParametricParams (excluding "selections")
			for k, v := range config.ParametricParams {
				if k != "selections" {
					catalogConfig.Parameters[k] = v
				}
			}
		}

		if selections, ok := config.ParametricParams["selections"]; ok {
			switch s := selections.(type) {
			case map[string]interface{}:
				for k, v := range s {
					if strVal, ok := v.(string); ok {
						catalogConfig.Selections[k] = strVal
					}
				}
			case map[string]string:
				for k, v := range s {
					catalogConfig.Selections[k] = v
				}
			}
		}
	}

	return catalogConfig
}
