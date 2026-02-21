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

type CategoryRepository struct {
	db *DB
}

func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, code, parent_id, name, description, image, sort_order, active,
		       pim_code, last_synced_at, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanCategory(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *CategoryRepository) GetByIDWithAncestors(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ancestors, err := r.GetAncestors(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Ancestors = ancestors
	return category, nil
}

func (r *CategoryRepository) GetAncestors(ctx context.Context, id uuid.UUID) ([]domain.Category, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			-- Start with the parent of the given category
			SELECT c.id, c.tenant_id, c.code, c.parent_id, c.name, c.description, c.image, 
			       c.sort_order, c.active, c.pim_code, c.last_synced_at, 
			       c.created_at, c.updated_at, c.deleted_at, 1 as depth
			FROM categories c
			INNER JOIN categories child ON child.parent_id = c.id
			WHERE child.id = $1 AND c.deleted_at IS NULL
			
			UNION ALL
			
			-- Recursively get parent categories
			SELECT c.id, c.tenant_id, c.code, c.parent_id, c.name, c.description, c.image,
			       c.sort_order, c.active, c.pim_code, c.last_synced_at,
			       c.created_at, c.updated_at, c.deleted_at, a.depth + 1
			FROM categories c
			INNER JOIN ancestors a ON a.parent_id = c.id
			WHERE c.deleted_at IS NULL
		)
		SELECT id, tenant_id, code, parent_id, name, description, image, sort_order, active,
		       pim_code, last_synced_at, created_at, updated_at, deleted_at
		FROM ancestors
		ORDER BY depth DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ancestors []domain.Category
	for rows.Next() {
		category, err := r.scanCategoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		ancestors = append(ancestors, *category)
	}

	return ancestors, rows.Err()
}

func (r *CategoryRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, code, parent_id, name, description, image, sort_order, active,
		       pim_code, last_synced_at, created_at, updated_at, deleted_at
		FROM categories
		WHERE tenant_id = $1 AND code = $2 AND deleted_at IS NULL
	`

	return r.scanCategory(r.db.Pool.QueryRow(ctx, query, tenantID, code))
}

func (r *CategoryRepository) GetTree(ctx context.Context, tenantID uuid.UUID) ([]domain.Category, error) {
	query := `
		WITH RECURSIVE category_tree AS (
			-- Root categories
			SELECT id, tenant_id, code, parent_id, name, description, image, sort_order, active,
			       pim_code, last_synced_at, created_at, updated_at, deleted_at, 0 as depth
			FROM categories
			WHERE tenant_id = $1 AND parent_id IS NULL AND deleted_at IS NULL
			
			UNION ALL
			
			-- Child categories
			SELECT c.id, c.tenant_id, c.code, c.parent_id, c.name, c.description, c.image,
			       c.sort_order, c.active, c.pim_code, c.last_synced_at, 
			       c.created_at, c.updated_at, c.deleted_at, ct.depth + 1
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
			WHERE c.deleted_at IS NULL
		),
		-- Get all child category IDs for each category (including self)
		category_descendants AS (
			SELECT id, id as descendant_id
			FROM category_tree
			
			UNION ALL
			
			SELECT cd.id, ct.id as descendant_id
			FROM category_descendants cd
			INNER JOIN category_tree ct ON ct.parent_id = cd.descendant_id
		),
		-- Count products for each category including descendants
		category_product_counts AS (
			SELECT cd.id as category_id,
			       COUNT(DISTINCT p.id)::int as product_count
			FROM category_descendants cd
			LEFT JOIN products p ON cd.descendant_id = ANY(p.category_ids) AND p.deleted_at IS NULL
			GROUP BY cd.id
		)
		SELECT ct.id, ct.tenant_id, ct.code, ct.parent_id, ct.name, ct.description, ct.image,
		       ct.sort_order, ct.active, ct.pim_code, ct.last_synced_at, 
		       ct.created_at, ct.updated_at, ct.deleted_at,
		       COALESCE(cpc.product_count, 0) as product_count
		FROM category_tree ct
		LEFT JOIN category_product_counts cpc ON ct.id = cpc.category_id
		ORDER BY depth, sort_order, code
	`

	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category, err := r.scanCategoryWithCountFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *category)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) List(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("parent_id = $%d", argNum))
		args = append(args, *filter.ParentID)
		argNum++
	}

	if filter.Active != nil {
		conditions = append(conditions, fmt.Sprintf("active = $%d", argNum))
		args = append(args, *filter.Active)
		argNum++
	}

	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf(`
			(code ILIKE $%d OR name::text ILIKE $%d)
		`, argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM categories WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query with product counts per category (including descendants)
	// Uses LATERAL join with recursive CTE per category for correct per-row counts
	query := fmt.Sprintf(`
		SELECT c.id, c.tenant_id, c.code, c.parent_id, c.name, c.description, c.image,
		       c.sort_order, c.active, c.pim_code, c.last_synced_at, 
		       c.created_at, c.updated_at, c.deleted_at,
		       COALESCE(pc.cnt, 0) as product_count
		FROM categories c
		LEFT JOIN LATERAL (
		    WITH RECURSIVE descendants AS (
		        SELECT c.id AS cat_id
		        UNION ALL
		        SELECT ch.id
		        FROM categories ch
		        INNER JOIN descendants d ON ch.parent_id = d.cat_id
		        WHERE ch.deleted_at IS NULL
		    )
		    SELECT COUNT(DISTINCT p.id)::int AS cnt
		    FROM products p
		    WHERE EXISTS (
		        SELECT 1 FROM descendants d WHERE d.cat_id = ANY(p.category_ids)
		    ) AND p.deleted_at IS NULL
		) pc ON true
		WHERE %s
		ORDER BY c.sort_order, c.code
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category, err := r.scanCategoryWithCountFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, *category)
	}

	return categories, total, rows.Err()
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	// Marshal JSON fields
	nameJSON, _ := json.Marshal(category.Name)
	descriptionJSON, _ := json.Marshal(category.Description)

	query := `
		INSERT INTO categories (
			id, tenant_id, code, parent_id, name, description, image, sort_order, active,
			pim_code, last_synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		category.ID,
		category.TenantID,
		category.Code,
		category.ParentID,
		nameJSON,
		descriptionJSON,
		category.Image,
		category.SortOrder,
		category.Active,
		category.PIMCode,
		category.LastSyncedAt,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	// Marshal JSON fields
	nameJSON, _ := json.Marshal(category.Name)
	descriptionJSON, _ := json.Marshal(category.Description)

	query := `
		UPDATE categories SET
			parent_id = $1,
			name = $2,
			description = $3,
			image = $4,
			sort_order = $5,
			active = $6,
			pim_code = $7,
			last_synced_at = $8,
			updated_at = $9
		WHERE id = $10 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query,
		category.ParentID,
		nameJSON,
		descriptionJSON,
		category.Image,
		category.SortOrder,
		category.Active,
		category.PIMCode,
		category.LastSyncedAt,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE categories SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) HasProducts(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM products 
			WHERE $1 = ANY(category_ids) AND deleted_at IS NULL
		)
	`

	var hasProducts bool
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&hasProducts)
	return hasProducts, err
}

func (r *CategoryRepository) scanCategory(row pgx.Row) (*domain.Category, error) {
	var category domain.Category
	var nameJSON, descriptionJSON []byte

	err := row.Scan(
		&category.ID,
		&category.TenantID,
		&category.Code,
		&category.ParentID,
		&nameJSON,
		&descriptionJSON,
		&category.Image,
		&category.SortOrder,
		&category.Active,
		&category.PIMCode,
		&category.LastSyncedAt,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	// Unmarshal JSON fields
	_ = json.Unmarshal(nameJSON, &category.Name)
	_ = json.Unmarshal(descriptionJSON, &category.Description)

	if category.Name == nil {
		category.Name = make(map[string]string)
	}
	if category.Description == nil {
		category.Description = make(map[string]string)
	}

	return &category, nil
}

func (r *CategoryRepository) scanCategoryFromRows(rows pgx.Rows) (*domain.Category, error) {
	var category domain.Category
	var nameJSON, descriptionJSON []byte

	err := rows.Scan(
		&category.ID,
		&category.TenantID,
		&category.Code,
		&category.ParentID,
		&nameJSON,
		&descriptionJSON,
		&category.Image,
		&category.SortOrder,
		&category.Active,
		&category.PIMCode,
		&category.LastSyncedAt,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	_ = json.Unmarshal(nameJSON, &category.Name)
	_ = json.Unmarshal(descriptionJSON, &category.Description)

	if category.Name == nil {
		category.Name = make(map[string]string)
	}
	if category.Description == nil {
		category.Description = make(map[string]string)
	}

	return &category, nil
}

func (r *CategoryRepository) scanCategoryWithCountFromRows(rows pgx.Rows) (*domain.Category, error) {
	var category domain.Category
	var nameJSON, descriptionJSON []byte

	err := rows.Scan(
		&category.ID,
		&category.TenantID,
		&category.Code,
		&category.ParentID,
		&nameJSON,
		&descriptionJSON,
		&category.Image,
		&category.SortOrder,
		&category.Active,
		&category.PIMCode,
		&category.LastSyncedAt,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
		&category.ProductCount,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	_ = json.Unmarshal(nameJSON, &category.Name)
	_ = json.Unmarshal(descriptionJSON, &category.Description)

	if category.Name == nil {
		category.Name = make(map[string]string)
	}
	if category.Description == nil {
		category.Description = make(map[string]string)
	}

	return &category, nil
}
