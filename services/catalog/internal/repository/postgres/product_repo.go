package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

type ProductRepository struct {
	db *DB
}

func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, sku, name, description, category_ids, attributes, status, images,
		       pim_identifier, last_synced_at, created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanProduct(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *ProductRepository) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, sku, name, description, category_ids, attributes, status, images,
		       pim_identifier, last_synced_at, created_at, updated_at, deleted_at
		FROM products
		WHERE tenant_id = $1 AND sku = $2 AND deleted_at IS NULL
	`

	return r.scanProduct(r.db.Pool.QueryRow(ctx, query, tenantID, sku))
}

func (r *ProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.CategoryID != nil {
		if filter.IncludeChildren {
			// Use recursive CTE to get all descendant categories
			conditions = append(conditions, fmt.Sprintf(`
				EXISTS (
					WITH RECURSIVE child_categories AS (
						SELECT id FROM categories WHERE id = $%d AND deleted_at IS NULL
						UNION ALL
						SELECT c.id FROM categories c
						INNER JOIN child_categories cc ON c.parent_id = cc.id
						WHERE c.deleted_at IS NULL
					)
					SELECT 1 FROM child_categories cc WHERE cc.id = ANY(category_ids)
				)
			`, argNum))
			args = append(args, *filter.CategoryID)
			argNum++
		} else {
			// Only direct category match
			conditions = append(conditions, fmt.Sprintf("$%d = ANY(category_ids)", argNum))
			args = append(args, *filter.CategoryID)
			argNum++
		}
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *filter.Status)
		argNum++
	}

	if len(filter.SKUs) > 0 {
		conditions = append(conditions, fmt.Sprintf("sku = ANY($%d)", argNum))
		args = append(args, filter.SKUs)
		argNum++
	}

	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf(`
			(sku ILIKE $%d OR name::text ILIKE $%d)
		`, argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, sku, name, description, category_ids, attributes, status, images,
		       pim_identifier, last_synced_at, created_at, updated_at, deleted_at
		FROM products
		WHERE %s
		ORDER BY sku
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, *product)
	}

	return products, total, rows.Err()
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	// Marshal JSON fields
	nameJSON, _ := json.Marshal(product.Name)
	descJSON, _ := json.Marshal(product.Description)
	attrJSON, _ := json.Marshal(product.Attributes)
	imagesJSON, _ := json.Marshal(product.Images)

	query := `
		INSERT INTO products (
			id, tenant_id, sku, name, description, category_ids, attributes, status, images,
			pim_identifier, last_synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		product.ID,
		product.TenantID,
		product.SKU,
		nameJSON,
		descJSON,
		product.CategoryIDs,
		attrJSON,
		product.Status,
		imagesJSON,
		product.PIMIdentifier,
		product.LastSyncedAt,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrProductAlreadyExists
		}
		return err
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	// Marshal JSON fields
	nameJSON, _ := json.Marshal(product.Name)
	descJSON, _ := json.Marshal(product.Description)
	attrJSON, _ := json.Marshal(product.Attributes)
	imagesJSON, _ := json.Marshal(product.Images)

	query := `
		UPDATE products SET
			name = $1,
			description = $2,
			category_ids = $3,
			attributes = $4,
			status = $5,
			images = $6,
			pim_identifier = $7,
			last_synced_at = $8,
			updated_at = $9
		WHERE id = $10 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query,
		nameJSON,
		descJSON,
		product.CategoryIDs,
		attrJSON,
		product.Status,
		imagesJSON,
		product.PIMIdentifier,
		product.LastSyncedAt,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE products SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) scanProduct(row pgx.Row) (*domain.Product, error) {
	var product domain.Product
	var nameJSON, descJSON, attrJSON, imagesJSON []byte

	err := row.Scan(
		&product.ID,
		&product.TenantID,
		&product.SKU,
		&nameJSON,
		&descJSON,
		&product.CategoryIDs,
		&attrJSON,
		&product.Status,
		&imagesJSON,
		&product.PIMIdentifier,
		&product.LastSyncedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	// Unmarshal JSON fields
	_ = json.Unmarshal(nameJSON, &product.Name)
	_ = json.Unmarshal(descJSON, &product.Description)
	_ = json.Unmarshal(attrJSON, &product.Attributes)
	_ = json.Unmarshal(imagesJSON, &product.Images)

	// Initialize empty slices if nil
	if product.CategoryIDs == nil {
		product.CategoryIDs = []uuid.UUID{}
	}
	if product.Attributes == nil {
		product.Attributes = []domain.ProductAttribute{}
	}
	if product.Images == nil {
		product.Images = []domain.ProductImage{}
	}
	if product.Name == nil {
		product.Name = make(map[string]string)
	}
	if product.Description == nil {
		product.Description = make(map[string]string)
	}

	return &product, nil
}

func (r *ProductRepository) scanProductFromRows(rows pgx.Rows) (*domain.Product, error) {
	var product domain.Product
	var nameJSON, descJSON, attrJSON, imagesJSON []byte

	err := rows.Scan(
		&product.ID,
		&product.TenantID,
		&product.SKU,
		&nameJSON,
		&descJSON,
		&product.CategoryIDs,
		&attrJSON,
		&product.Status,
		&imagesJSON,
		&product.PIMIdentifier,
		&product.LastSyncedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	_ = json.Unmarshal(nameJSON, &product.Name)
	_ = json.Unmarshal(descJSON, &product.Description)
	_ = json.Unmarshal(attrJSON, &product.Attributes)
	_ = json.Unmarshal(imagesJSON, &product.Images)

	// Initialize empty slices if nil
	if product.CategoryIDs == nil {
		product.CategoryIDs = []uuid.UUID{}
	}
	if product.Attributes == nil {
		product.Attributes = []domain.ProductAttribute{}
	}
	if product.Images == nil {
		product.Images = []domain.ProductImage{}
	}
	if product.Name == nil {
		product.Name = make(map[string]string)
	}
	if product.Description == nil {
		product.Description = make(map[string]string)
	}

	return &product, nil
}
