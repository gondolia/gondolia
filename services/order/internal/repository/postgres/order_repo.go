package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/order/internal/domain"
)

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount, total,
		       currency, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order, err := r.scanOrder(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	// Load items
	items, err := r.GetItemsByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	// Load status history
	history, err := r.GetStatusHistory(ctx, id)
	if err != nil {
		return nil, err
	}
	order.StatusHistory = history

	return order, nil
}

func (r *OrderRepository) GetByOrderNumber(ctx context.Context, tenantID uuid.UUID, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount, total,
		       currency, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders
		WHERE tenant_id = $1 AND order_number = $2
	`

	order, err := r.scanOrder(r.db.Pool.QueryRow(ctx, query, tenantID, orderNumber))
	if err != nil {
		return nil, err
	}

	// Load items
	items, err := r.GetItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	// Load status history
	history, err := r.GetStatusHistory(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.StatusHistory = history

	return order, nil
}

func (r *OrderRepository) List(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argNum))
		args = append(args, *filter.UserID)
		argNum++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *filter.Status)
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query with item count subquery
	query := fmt.Sprintf(`
		SELECT o.id, o.tenant_id, o.user_id, o.order_number, o.status, o.subtotal, o.tax_amount, o.total,
		       o.currency, o.shipping_address, o.billing_address, o.notes, o.created_at, o.updated_at,
		       (SELECT COUNT(*) FROM order_items WHERE order_id = o.id) as item_count
		FROM orders o
		WHERE %s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order, err := r.scanOrderFromRowsWithCount(rows)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, *order)
	}

	return orders, total, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	shippingAddrJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	billingAddrJSON, err := json.Marshal(order.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal billing address: %w", err)
	}

	query := `
		INSERT INTO orders (id, tenant_id, user_id, order_number, status, subtotal, tax_amount, total,
		                    currency, shipping_address, billing_address, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err = r.db.Pool.Exec(ctx, query,
		order.ID,
		order.TenantID,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.Subtotal,
		order.TaxAmount,
		order.Total,
		order.Currency,
		shippingAddrJSON,
		billingAddrJSON,
		order.Notes,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	shippingAddrJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	billingAddrJSON, err := json.Marshal(order.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal billing address: %w", err)
	}

	query := `
		UPDATE orders
		SET status = $2, subtotal = $3, tax_amount = $4, total = $5, currency = $6,
		    shipping_address = $7, billing_address = $8, notes = $9, updated_at = $10
		WHERE id = $1
	`

	_, err = r.db.Pool.Exec(ctx, query,
		order.ID,
		order.Status,
		order.Subtotal,
		order.TaxAmount,
		order.Total,
		order.Currency,
		shippingAddrJSON,
		billingAddrJSON,
		order.Notes,
		order.UpdatedAt,
	)

	return err
}

func (r *OrderRepository) CreateItems(ctx context.Context, items []domain.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	// Build bulk insert query
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]any, 0, len(items)*12)
	argNum := 1

	for _, item := range items {
		configJSON, err := json.Marshal(item.Configuration)
		if err != nil {
			return fmt.Errorf("failed to marshal configuration: %w", err)
		}

		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argNum, argNum+1, argNum+2, argNum+3, argNum+4, argNum+5,
			argNum+6, argNum+7, argNum+8, argNum+9, argNum+10, argNum+11,
		))

		valueArgs = append(valueArgs,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.VariantID,
			item.ProductType,
			item.ProductName,
			item.SKU,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
			item.Currency,
			configJSON,
		)

		argNum += 12
	}

	query := fmt.Sprintf(`
		INSERT INTO order_items (id, order_id, product_id, variant_id, product_type, product_name,
		                         sku, quantity, unit_price, total_price, currency, configuration)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.Pool.Exec(ctx, query, valueArgs...)
	return err
}

func (r *OrderRepository) GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, variant_id, product_type, product_name, sku,
		       quantity, unit_price, total_price, currency, configuration, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		item, err := r.scanOrderItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

func (r *OrderRepository) AddStatusLog(ctx context.Context, log *domain.OrderStatusLog) error {
	query := `
		INSERT INTO order_status_history (id, order_id, from_status, to_status, changed_by, note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		log.ID,
		log.OrderID,
		log.FromStatus,
		log.ToStatus,
		log.ChangedBy,
		log.Note,
		log.CreatedAt,
	)

	return err
}

func (r *OrderRepository) GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]domain.OrderStatusLog, error) {
	query := `
		SELECT id, order_id, from_status, to_status, changed_by, note, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.OrderStatusLog
	for rows.Next() {
		log, err := r.scanStatusLog(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, *log)
	}

	return logs, nil
}

func (r *OrderRepository) GenerateOrderNumber(ctx context.Context, tenantID uuid.UUID) (string, error) {
	// Order number format: ORD-{YYYYMMDD}-{XXXX} (sequentiell pro Tag)
	now := time.Now()
	datePrefix := now.Format("20060102")

	// Get the next sequence number for today
	query := `
		SELECT COUNT(*) + 1
		FROM orders
		WHERE tenant_id = $1 AND order_number LIKE $2
	`

	var sequence int
	err := r.db.Pool.QueryRow(ctx, query, tenantID, fmt.Sprintf("ORD-%s-%%", datePrefix)).Scan(&sequence)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("ORD-%s-%04d", datePrefix, sequence), nil
}

func (r *OrderRepository) scanOrder(row pgx.Row) (*domain.Order, error) {
	var order domain.Order
	var shippingAddrJSON, billingAddrJSON []byte

	err := row.Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.Total,
		&order.Currency,
		&shippingAddrJSON,
		&billingAddrJSON,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
	}

	if err := json.Unmarshal(billingAddrJSON, &order.BillingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) scanOrderFromRows(rows pgx.Rows) (*domain.Order, error) {
	var order domain.Order
	var shippingAddrJSON, billingAddrJSON []byte

	err := rows.Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.Total,
		&order.Currency,
		&shippingAddrJSON,
		&billingAddrJSON,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
	}

	if err := json.Unmarshal(billingAddrJSON, &order.BillingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) scanOrderFromRowsWithCount(rows pgx.Rows) (*domain.Order, error) {
	var order domain.Order
	var shippingAddrJSON, billingAddrJSON []byte

	err := rows.Scan(
		&order.ID,
		&order.TenantID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.Total,
		&order.Currency,
		&shippingAddrJSON,
		&billingAddrJSON,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.ItemCount,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
	}

	if err := json.Unmarshal(billingAddrJSON, &order.BillingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) scanOrderItem(rows pgx.Rows) (*domain.OrderItem, error) {
	var item domain.OrderItem
	var configJSON []byte

	err := rows.Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.VariantID,
		&item.ProductType,
		&item.ProductName,
		&item.SKU,
		&item.Quantity,
		&item.UnitPrice,
		&item.TotalPrice,
		&item.Currency,
		&configJSON,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if configJSON != nil && len(configJSON) > 0 {
		if err := json.Unmarshal(configJSON, &item.Configuration); err != nil {
			return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
		}
	}

	return &item, nil
}

func (r *OrderRepository) scanStatusLog(rows pgx.Rows) (*domain.OrderStatusLog, error) {
	var log domain.OrderStatusLog
	var fromStatus *string

	err := rows.Scan(
		&log.ID,
		&log.OrderID,
		&fromStatus,
		&log.ToStatus,
		&log.ChangedBy,
		&log.Note,
		&log.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if fromStatus != nil {
		status := domain.OrderStatus(*fromStatus)
		log.FromStatus = &status
	}

	return &log, nil
}
