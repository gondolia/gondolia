package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

type ParametricPricingRepository struct {
	db *DB
}

func NewParametricPricingRepository(db *DB) *ParametricPricingRepository {
	return &ParametricPricingRepository{db: db}
}

func (r *ParametricPricingRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.ParametricPricing, error) {
	var pp domain.ParametricPricing
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, product_id, formula_type, base_price, unit_price, currency, min_order_value, created_at, updated_at
		 FROM parametric_pricing WHERE product_id = $1`, productID).Scan(
		&pp.ID, &pp.ProductID, &pp.FormulaType, &pp.BasePrice, &pp.UnitPrice, &pp.Currency, &pp.MinOrderValue, &pp.CreatedAt, &pp.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrParametricPricingNotFound
		}
		return nil, err
	}
	return &pp, nil
}

func (r *ParametricPricingRepository) Create(ctx context.Context, pricing *domain.ParametricPricing) error {
	pricing.ID = uuid.New()
	pricing.CreatedAt = time.Now()
	pricing.UpdatedAt = pricing.CreatedAt

	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO parametric_pricing (id, product_id, formula_type, base_price, unit_price, currency, min_order_value, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		pricing.ID, pricing.ProductID, pricing.FormulaType, pricing.BasePrice, pricing.UnitPrice, pricing.Currency, pricing.MinOrderValue, pricing.CreatedAt, pricing.UpdatedAt)
	return err
}

func (r *ParametricPricingRepository) Update(ctx context.Context, pricing *domain.ParametricPricing) error {
	pricing.UpdatedAt = time.Now()
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE parametric_pricing SET formula_type=$1, base_price=$2, unit_price=$3, currency=$4, min_order_value=$5, updated_at=$6
		 WHERE product_id = $7`,
		pricing.FormulaType, pricing.BasePrice, pricing.UnitPrice, pricing.Currency, pricing.MinOrderValue, pricing.UpdatedAt, pricing.ProductID)
	return err
}

func (r *ParametricPricingRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM parametric_pricing WHERE product_id = $1`, productID)
	return err
}

// AxisOptionRepository handles axis_options for parametric products
type AxisOptionRepository struct {
	db *DB
}

func NewAxisOptionRepository(db *DB) *AxisOptionRepository {
	return &AxisOptionRepository{db: db}
}

func (r *AxisOptionRepository) ListByAxisID(ctx context.Context, axisID uuid.UUID) ([]domain.AxisOption, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT code, label, position FROM axis_options WHERE axis_id = $1 ORDER BY position`, axisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []domain.AxisOption
	for rows.Next() {
		var opt domain.AxisOption
		var labelJSON []byte
		if err := rows.Scan(&opt.Code, &labelJSON, &opt.Position); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(labelJSON, &opt.Label)
		options = append(options, opt)
	}
	return options, rows.Err()
}

// SKUMappingRepository handles parametric_sku_mapping queries
type SKUMappingRepository struct {
	db *DB
}

func NewSKUMappingRepository(db *DB) *SKUMappingRepository {
	return &SKUMappingRepository{db: db}
}

func (r *SKUMappingRepository) FindBySelections(ctx context.Context, productID uuid.UUID, selections map[string]string) (*domain.SKUMapping, error) {
	selectionsJSON, err := json.Marshal(selections)
	if err != nil {
		return nil, err
	}

	var m domain.SKUMapping
	var selJSON []byte
	err = r.db.Pool.QueryRow(ctx,
		`SELECT id, product_id, selections, sku, unit_price, base_price, stock, created_at, updated_at
		 FROM parametric_sku_mapping
		 WHERE product_id = $1 AND selections @> $2::jsonb AND $2::jsonb @> selections`,
		productID, selectionsJSON).Scan(
		&m.ID, &m.ProductID, &selJSON, &m.SKU, &m.UnitPrice, &m.BasePrice, &m.Stock, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSKUMappingNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(selJSON, &m.Selections)
	return &m, nil
}

func (r *SKUMappingRepository) ListByProductID(ctx context.Context, productID uuid.UUID) ([]domain.SKUMapping, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, product_id, selections, sku, unit_price, base_price, stock, created_at, updated_at
		 FROM parametric_sku_mapping WHERE product_id = $1 ORDER BY sku`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []domain.SKUMapping
	for rows.Next() {
		var m domain.SKUMapping
		var selJSON []byte
		if err := rows.Scan(&m.ID, &m.ProductID, &selJSON, &m.SKU, &m.UnitPrice, &m.BasePrice, &m.Stock, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(selJSON, &m.Selections)
		mappings = append(mappings, m)
	}
	return mappings, rows.Err()
}

func (r *AxisOptionRepository) ListByProductID(ctx context.Context, productID uuid.UUID) (map[uuid.UUID][]domain.AxisOption, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT ao.axis_id, ao.code, ao.label, ao.position
		 FROM axis_options ao
		 JOIN variant_axes va ON va.id = ao.axis_id
		 WHERE va.product_id = $1
		 ORDER BY ao.axis_id, ao.position`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]domain.AxisOption)
	for rows.Next() {
		var axisID uuid.UUID
		var opt domain.AxisOption
		var labelJSON []byte
		if err := rows.Scan(&axisID, &opt.Code, &labelJSON, &opt.Position); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(labelJSON, &opt.Label)
		result[axisID] = append(result[axisID], opt)
	}
	return result, rows.Err()
}
