package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/cart/internal/domain"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetByCode(ctx context.Context, code string) (*domain.Tenant, error)
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Cart, error)
	GetActiveCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) (*domain.Cart, error)
	Create(ctx context.Context, cart *domain.Cart) error
	Update(ctx context.Context, cart *domain.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Cart items
	AddItem(ctx context.Context, item *domain.CartItem) error
	UpdateItem(ctx context.Context, item *domain.CartItem) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
	GetItem(ctx context.Context, id uuid.UUID) (*domain.CartItem, error)
	GetCartItems(ctx context.Context, cartID uuid.UUID) ([]domain.CartItem, error)
	FindMatchingItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, configHash string) (*domain.CartItem, error)
	ClearCart(ctx context.Context, cartID uuid.UUID) error

	// Merge carts (guest -> user)
	MergeCarts(ctx context.Context, fromCartID, toCartID uuid.UUID) error
}
