package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/cart/internal/domain"
)

type CartRepository struct {
	db *DB
}

func NewCartRepository(db *DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cart, error) {
	query := `
		SELECT id, tenant_id, user_id, session_id, status, created_at, updated_at
		FROM carts
		WHERE id = $1
	`

	cart, err := r.scanCart(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, err
	}

	// Load items
	items, err := r.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	cart.Items = items

	// Compute subtotal and currency from items
	cart.Subtotal = cart.TotalPrice()
	cart.Currency = cart.GetCurrency()

	return cart, nil
}

func (r *CartRepository) GetActiveCart(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) (*domain.Cart, error) {
	var cart *domain.Cart
	var err error

	// Strategy: Try to find cart by user_id first, then by session_id
	// If a session cart is found and user_id is provided, claim the cart for the user

	if userID != nil {
		// First, try to find a cart by user_id
		query := `
			SELECT id, tenant_id, user_id, session_id, status, created_at, updated_at
			FROM carts
			WHERE tenant_id = $1 AND user_id = $2 AND status = $3
			ORDER BY updated_at DESC
			LIMIT 1
		`
		cart, err = r.scanCart(r.db.Pool.QueryRow(ctx, query, tenantID, userID, domain.CartStatusActive))
		if err == nil {
			// Found user cart, load items and return
			items, err := r.GetCartItems(ctx, cart.ID)
			if err != nil {
				return nil, err
			}
			cart.Items = items
			cart.Subtotal = cart.TotalPrice()
			cart.Currency = cart.GetCurrency()
			return cart, nil
		}

		// No user cart found, try session_id if provided
		if sessionID != nil && *sessionID != "" {
			query = `
				SELECT id, tenant_id, user_id, session_id, status, created_at, updated_at
				FROM carts
				WHERE tenant_id = $1 AND session_id = $2 AND status = $3
				ORDER BY updated_at DESC
				LIMIT 1
			`
			cart, err = r.scanCart(r.db.Pool.QueryRow(ctx, query, tenantID, sessionID, domain.CartStatusActive))
			if err == nil {
				// Found session cart - claim it for the user
				cart.UserID = userID
				updateQuery := `UPDATE carts SET user_id = $1, updated_at = NOW() WHERE id = $2`
				if _, err := r.db.Pool.Exec(ctx, updateQuery, userID, cart.ID); err != nil {
					return nil, err
				}

				// Load items
				items, err := r.GetCartItems(ctx, cart.ID)
				if err != nil {
					return nil, err
				}
				cart.Items = items
				cart.Subtotal = cart.TotalPrice()
				cart.Currency = cart.GetCurrency()
				return cart, nil
			}
		}

		// No cart found at all
		return nil, domain.ErrCartNotFound
	}

	// No user_id, try session_id only (guest flow)
	if sessionID != nil && *sessionID != "" {
		query := `
			SELECT id, tenant_id, user_id, session_id, status, created_at, updated_at
			FROM carts
			WHERE tenant_id = $1 AND session_id = $2 AND status = $3
			ORDER BY updated_at DESC
			LIMIT 1
		`
		cart, err = r.scanCart(r.db.Pool.QueryRow(ctx, query, tenantID, sessionID, domain.CartStatusActive))
		if err != nil {
			return nil, err
		}

		// Load items
		items, err := r.GetCartItems(ctx, cart.ID)
		if err != nil {
			return nil, err
		}
		cart.Items = items
		cart.Subtotal = cart.TotalPrice()
		cart.Currency = cart.GetCurrency()
		return cart, nil
	}

	return nil, domain.ErrCartNotFound
}

func (r *CartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	query := `
		INSERT INTO carts (id, tenant_id, user_id, session_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		cart.ID,
		cart.TenantID,
		cart.UserID,
		cart.SessionID,
		cart.Status,
		cart.CreatedAt,
		cart.UpdatedAt,
	)
	return err
}

func (r *CartRepository) Update(ctx context.Context, cart *domain.Cart) error {
	query := `
		UPDATE carts
		SET user_id = $2, session_id = $3, status = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query,
		cart.ID,
		cart.UserID,
		cart.SessionID,
		cart.Status,
		cart.UpdatedAt,
	)
	return err
}

func (r *CartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM carts WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *CartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	var configJSON []byte
	var err error
	if item.Configuration != nil {
		configJSON, err = json.Marshal(item.Configuration)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO cart_items (id, cart_id, product_id, variant_id, product_type, product_name, sku, image_url, quantity, unit_price, total_price, currency, configuration, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = r.db.Pool.Exec(ctx, query,
		item.ID,
		item.CartID,
		item.ProductID,
		item.VariantID,
		item.ProductType,
		item.ProductName,
		item.SKU,
		item.ImageURL,
		item.Quantity,
		item.UnitPrice,
		item.TotalPrice,
		item.Currency,
		configJSON,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

func (r *CartRepository) UpdateItem(ctx context.Context, item *domain.CartItem) error {
	var configJSON []byte
	var err error
	if item.Configuration != nil {
		configJSON, err = json.Marshal(item.Configuration)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE cart_items
		SET quantity = $2, unit_price = $3, total_price = $4, currency = $5, product_name = $6, sku = $7, image_url = $8, configuration = $9, updated_at = $10
		WHERE id = $1
	`

	_, err = r.db.Pool.Exec(ctx, query,
		item.ID,
		item.Quantity,
		item.UnitPrice,
		item.TotalPrice,
		item.Currency,
		item.ProductName,
		item.SKU,
		item.ImageURL,
		configJSON,
		item.UpdatedAt,
	)
	return err
}

func (r *CartRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *CartRepository) GetItem(ctx context.Context, id uuid.UUID) (*domain.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, variant_id, product_type, product_name, sku, image_url, quantity, unit_price, total_price, currency, configuration, created_at, updated_at
		FROM cart_items
		WHERE id = $1
	`

	return r.scanCartItem(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *CartRepository) GetCartItems(ctx context.Context, cartID uuid.UUID) ([]domain.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, variant_id, product_type, product_name, sku, image_url, quantity, unit_price, total_price, currency, configuration, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		item, err := r.scanCartItemFromRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (r *CartRepository) FindMatchingItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, configHash string) (*domain.CartItem, error) {
	var query string
	var args []interface{}

	if variantID != nil {
		query = `
			SELECT id, cart_id, product_id, variant_id, product_type, product_name, sku, image_url, quantity, unit_price, total_price, currency, configuration, created_at, updated_at
			FROM cart_items
			WHERE cart_id = $1 AND product_id = $2 AND variant_id = $3 AND COALESCE(md5(configuration::text), '') = $4
			LIMIT 1
		`
		args = []interface{}{cartID, productID, variantID, configHash}
	} else {
		query = `
			SELECT id, cart_id, product_id, variant_id, product_type, product_name, sku, image_url, quantity, unit_price, total_price, currency, configuration, created_at, updated_at
			FROM cart_items
			WHERE cart_id = $1 AND product_id = $2 AND variant_id IS NULL AND COALESCE(md5(configuration::text), '') = $4
			LIMIT 1
		`
		args = []interface{}{cartID, productID, nil, configHash}
	}

	return r.scanCartItem(r.db.Pool.QueryRow(ctx, query, args...))
}

func (r *CartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, cartID)
	return err
}

func (r *CartRepository) MergeCarts(ctx context.Context, fromCartID, toCartID uuid.UUID) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Move items from fromCart to toCart
	query := `
		UPDATE cart_items
		SET cart_id = $1
		WHERE cart_id = $2
	`
	_, err = tx.Exec(ctx, query, toCartID, fromCartID)
	if err != nil {
		return err
	}

	// Mark fromCart as merged
	query = `
		UPDATE carts
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, query, domain.CartStatusMerged, fromCartID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *CartRepository) scanCart(row pgx.Row) (*domain.Cart, error) {
	var cart domain.Cart
	err := row.Scan(
		&cart.ID,
		&cart.TenantID,
		&cart.UserID,
		&cart.SessionID,
		&cart.Status,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) scanCartItem(row pgx.Row) (*domain.CartItem, error) {
	var item domain.CartItem
	var configJSON []byte

	err := row.Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.VariantID,
		&item.ProductType,
		&item.ProductName,
		&item.SKU,
		&item.ImageURL,
		&item.Quantity,
		&item.UnitPrice,
		&item.TotalPrice,
		&item.Currency,
		&configJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCartItemNotFound
		}
		return nil, err
	}

	// Unmarshal configuration
	if len(configJSON) > 0 {
		var config domain.Configuration
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, err
		}
		item.Configuration = &config
	}

	return &item, nil
}

func (r *CartRepository) scanCartItemFromRows(rows pgx.Rows) (*domain.CartItem, error) {
	var item domain.CartItem
	var configJSON []byte

	err := rows.Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.VariantID,
		&item.ProductType,
		&item.ProductName,
		&item.SKU,
		&item.ImageURL,
		&item.Quantity,
		&item.UnitPrice,
		&item.TotalPrice,
		&item.Currency,
		&configJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal configuration
	if len(configJSON) > 0 {
		var config domain.Configuration
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, err
		}
		item.Configuration = &config
	}

	return &item, nil
}
