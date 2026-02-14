package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

type PriceRepository struct {
	db *DB
}

func NewPriceRepository(db *DB) *PriceRepository {
	return &PriceRepository{db: db}
}

func (r *PriceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	query := `
		SELECT id, tenant_id, product_id, customer_group_id, min_quantity, price, currency,
		       valid_from, valid_to, created_at, updated_at, deleted_at
		FROM prices
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanPrice(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *PriceRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Price, error) {
	query := `
		SELECT id, tenant_id, product_id, customer_group_id, min_quantity, price, currency,
		       valid_from, valid_to, created_at, updated_at, deleted_at
		FROM prices
		WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY customer_group_id NULLS FIRST, min_quantity
	`

	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []domain.Price
	for rows.Next() {
		price, err := r.scanPriceFromRows(rows)
		if err != nil {
			return nil, err
		}
		prices = append(prices, *price)
	}

	return prices, rows.Err()
}

func (r *PriceRepository) List(ctx context.Context, filter domain.PriceFilter) ([]domain.Price, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", argNum))
		args = append(args, *filter.ProductID)
		argNum++
	}

	if filter.CustomerGroupID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_group_id = $%d", argNum))
		args = append(args, *filter.CustomerGroupID)
		argNum++
	}

	if filter.Currency != nil {
		conditions = append(conditions, fmt.Sprintf("currency = $%d", argNum))
		args = append(args, *filter.Currency)
		argNum++
	}

	if filter.ValidAt != nil {
		conditions = append(conditions, fmt.Sprintf(`
			(valid_from IS NULL OR valid_from <= $%d) AND
			(valid_to IS NULL OR valid_to >= $%d)
		`, argNum, argNum))
		args = append(args, *filter.ValidAt)
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM prices WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, product_id, customer_group_id, min_quantity, price, currency,
		       valid_from, valid_to, created_at, updated_at, deleted_at
		FROM prices
		WHERE %s
		ORDER BY product_id, customer_group_id NULLS FIRST, min_quantity
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var prices []domain.Price
	for rows.Next() {
		price, err := r.scanPriceFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		prices = append(prices, *price)
	}

	return prices, total, rows.Err()
}

func (r *PriceRepository) Create(ctx context.Context, price *domain.Price) error {
	query := `
		INSERT INTO prices (
			id, tenant_id, product_id, customer_group_id, min_quantity, price, currency,
			valid_from, valid_to, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		price.ID,
		price.TenantID,
		price.ProductID,
		price.CustomerGroupID,
		price.MinQuantity,
		price.Price,
		price.Currency,
		price.ValidFrom,
		price.ValidTo,
		price.CreatedAt,
		price.UpdatedAt,
	)

	return err
}

func (r *PriceRepository) Update(ctx context.Context, price *domain.Price) error {
	query := `
		UPDATE prices SET
			customer_group_id = $1,
			min_quantity = $2,
			price = $3,
			currency = $4,
			valid_from = $5,
			valid_to = $6,
			updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query,
		price.CustomerGroupID,
		price.MinQuantity,
		price.Price,
		price.Currency,
		price.ValidFrom,
		price.ValidTo,
		price.UpdatedAt,
		price.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPriceNotFound
	}

	return nil
}

func (r *PriceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE prices SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPriceNotFound
	}

	return nil
}

func (r *PriceRepository) CheckOverlap(ctx context.Context, price *domain.Price) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM prices 
			WHERE tenant_id = $1 
			  AND product_id = $2 
			  AND (customer_group_id = $3 OR (customer_group_id IS NULL AND $3 IS NULL))
			  AND min_quantity = $4
			  AND id != $5
			  AND deleted_at IS NULL
			  AND (
				  (valid_from IS NULL AND valid_to IS NULL) OR
				  (valid_from IS NULL AND valid_to >= $6) OR
				  (valid_to IS NULL AND valid_from <= $7) OR
				  (valid_from <= $7 AND valid_to >= $6)
			  )
		)
	`

	var hasOverlap bool
	err := r.db.Pool.QueryRow(ctx, query,
		price.TenantID,
		price.ProductID,
		price.CustomerGroupID,
		price.MinQuantity,
		price.ID,
		price.ValidFrom,
		price.ValidTo,
	).Scan(&hasOverlap)

	return hasOverlap, err
}

func (r *PriceRepository) scanPrice(row pgx.Row) (*domain.Price, error) {
	var price domain.Price

	err := row.Scan(
		&price.ID,
		&price.TenantID,
		&price.ProductID,
		&price.CustomerGroupID,
		&price.MinQuantity,
		&price.Price,
		&price.Currency,
		&price.ValidFrom,
		&price.ValidTo,
		&price.CreatedAt,
		&price.UpdatedAt,
		&price.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPriceNotFound
		}
		return nil, err
	}

	return &price, nil
}

func (r *PriceRepository) scanPriceFromRows(rows pgx.Rows) (*domain.Price, error) {
	var price domain.Price

	err := rows.Scan(
		&price.ID,
		&price.TenantID,
		&price.ProductID,
		&price.CustomerGroupID,
		&price.MinQuantity,
		&price.Price,
		&price.Currency,
		&price.ValidFrom,
		&price.ValidTo,
		&price.CreatedAt,
		&price.UpdatedAt,
		&price.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &price, nil
}
