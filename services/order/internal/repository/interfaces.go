package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/order/internal/domain"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetByCode(ctx context.Context, code string) (*domain.Tenant, error)
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, tenantID uuid.UUID, orderNumber string) (*domain.Order, error)
	List(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int, error)
	Create(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error

	// Order items
	CreateItems(ctx context.Context, items []domain.OrderItem) error
	GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error)

	// Status history
	AddStatusLog(ctx context.Context, log *domain.OrderStatusLog) error
	GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]domain.OrderStatusLog, error)

	// Order number generation
	GenerateOrderNumber(ctx context.Context, tenantID uuid.UUID) (string, error)
}
