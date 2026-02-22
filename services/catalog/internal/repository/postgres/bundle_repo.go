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

type BundleRepository struct {
	db *DB
}

func NewBundleRepository(db *DB) *BundleRepository {
	return &BundleRepository{db: db}
}

// GetComponents retrieves all components for a bundle product
func (r *BundleRepository) GetComponents(ctx context.Context, bundleProductID uuid.UUID) ([]domain.BundleComponent, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, tenant_id, bundle_product_id, component_product_id, quantity,
		        min_quantity, max_quantity, sort_order, default_parameters,
		        created_at, updated_at
		 FROM bundle_components
		 WHERE bundle_product_id = $1
		 ORDER BY sort_order, created_at`, bundleProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []domain.BundleComponent
	for rows.Next() {
		var c domain.BundleComponent
		var defaultParamsJSON []byte

		err := rows.Scan(
			&c.ID, &c.TenantID, &c.BundleProductID, &c.ComponentProductID,
			&c.Quantity, &c.MinQuantity, &c.MaxQuantity, &c.SortOrder,
			&defaultParamsJSON, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse default_parameters JSON if present
		if len(defaultParamsJSON) > 0 {
			if err := json.Unmarshal(defaultParamsJSON, &c.DefaultParameters); err != nil {
				return nil, err
			}
		}

		components = append(components, c)
	}

	return components, rows.Err()
}

// GetComponentByID retrieves a single component by its ID
func (r *BundleRepository) GetComponentByID(ctx context.Context, componentID uuid.UUID) (*domain.BundleComponent, error) {
	var c domain.BundleComponent
	var defaultParamsJSON []byte

	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, bundle_product_id, component_product_id, quantity,
		        min_quantity, max_quantity, sort_order, default_parameters,
		        created_at, updated_at
		 FROM bundle_components
		 WHERE id = $1`, componentID).Scan(
		&c.ID, &c.TenantID, &c.BundleProductID, &c.ComponentProductID,
		&c.Quantity, &c.MinQuantity, &c.MaxQuantity, &c.SortOrder,
		&defaultParamsJSON, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrComponentNotFound
		}
		return nil, err
	}

	// Parse default_parameters JSON if present
	if len(defaultParamsJSON) > 0 {
		if err := json.Unmarshal(defaultParamsJSON, &c.DefaultParameters); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// SetComponents replaces all components for a bundle (delete + insert)
func (r *BundleRepository) SetComponents(ctx context.Context, bundleProductID uuid.UUID, components []domain.BundleComponent) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing components
	_, err = tx.Exec(ctx, `DELETE FROM bundle_components WHERE bundle_product_id = $1`, bundleProductID)
	if err != nil {
		return err
	}

	// Insert new components
	for _, c := range components {
		c.ID = uuid.New()
		c.BundleProductID = bundleProductID
		c.CreatedAt = time.Now()
		c.UpdatedAt = c.CreatedAt

		// Marshal default_parameters to JSON
		var defaultParamsJSON []byte
		if c.DefaultParameters != nil {
			defaultParamsJSON, err = json.Marshal(c.DefaultParameters)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO bundle_components (
				id, tenant_id, bundle_product_id, component_product_id,
				quantity, min_quantity, max_quantity, sort_order,
				default_parameters, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			c.ID, c.TenantID, bundleProductID, c.ComponentProductID,
			c.Quantity, c.MinQuantity, c.MaxQuantity, c.SortOrder,
			defaultParamsJSON, c.CreatedAt, c.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
