package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

type AttributeTranslationRepository struct {
	db *DB
}

func NewAttributeTranslationRepository(db *DB) *AttributeTranslationRepository {
	return &AttributeTranslationRepository{db: db}
}

func (r *AttributeTranslationRepository) GetByKey(ctx context.Context, tenantID uuid.UUID, attributeKey, locale string) (*domain.AttributeTranslation, error) {
	query := `
		SELECT id, tenant_id, attribute_key, locale, display_name, unit, description, created_at, updated_at
		FROM attribute_translations
		WHERE tenant_id = $1 AND attribute_key = $2 AND locale = $3
	`

	return r.scanTranslation(r.db.Pool.QueryRow(ctx, query, tenantID, attributeKey, locale))
}

func (r *AttributeTranslationRepository) GetByTenantAndLocale(ctx context.Context, tenantID uuid.UUID, locale string) (map[string]*domain.AttributeTranslation, error) {
	query := `
		SELECT id, tenant_id, attribute_key, locale, display_name, unit, description, created_at, updated_at
		FROM attribute_translations
		WHERE tenant_id = $1 AND locale = $2
		ORDER BY attribute_key
	`

	rows, err := r.db.Pool.Query(ctx, query, tenantID, locale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	translations := make(map[string]*domain.AttributeTranslation)
	for rows.Next() {
		translation, err := r.scanTranslationFromRows(rows)
		if err != nil {
			return nil, err
		}
		translations[translation.AttributeKey] = translation
	}

	return translations, rows.Err()
}

func (r *AttributeTranslationRepository) List(ctx context.Context, filter domain.AttributeTranslationFilter) ([]domain.AttributeTranslation, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	if filter.AttributeKey != nil {
		conditions = append(conditions, fmt.Sprintf("attribute_key = $%d", argNum))
		args = append(args, *filter.AttributeKey)
		argNum++
	}

	if filter.Locale != nil {
		conditions = append(conditions, fmt.Sprintf("locale = $%d", argNum))
		args = append(args, *filter.Locale)
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM attribute_translations WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, attribute_key, locale, display_name, unit, description, created_at, updated_at
		FROM attribute_translations
		WHERE %s
		ORDER BY attribute_key, locale
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var translations []domain.AttributeTranslation
	for rows.Next() {
		translation, err := r.scanTranslationFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		translations = append(translations, *translation)
	}

	return translations, total, rows.Err()
}

func (r *AttributeTranslationRepository) Create(ctx context.Context, translation *domain.AttributeTranslation) error {
	query := `
		INSERT INTO attribute_translations (
			id, tenant_id, attribute_key, locale, display_name, unit, description, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		translation.ID,
		translation.TenantID,
		translation.AttributeKey,
		translation.Locale,
		translation.DisplayName,
		translation.Unit,
		translation.Description,
		translation.CreatedAt,
		translation.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("translation for attribute '%s' in locale '%s' already exists", translation.AttributeKey, translation.Locale)
		}
		return err
	}

	return nil
}

func (r *AttributeTranslationRepository) Update(ctx context.Context, translation *domain.AttributeTranslation) error {
	query := `
		UPDATE attribute_translations SET
			display_name = $1,
			unit = $2,
			description = $3,
			updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.Pool.Exec(ctx, query,
		translation.DisplayName,
		translation.Unit,
		translation.Description,
		translation.UpdatedAt,
		translation.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attribute translation not found")
	}

	return nil
}

func (r *AttributeTranslationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM attribute_translations WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attribute translation not found")
	}

	return nil
}

func (r *AttributeTranslationRepository) scanTranslation(row pgx.Row) (*domain.AttributeTranslation, error) {
	var translation domain.AttributeTranslation

	err := row.Scan(
		&translation.ID,
		&translation.TenantID,
		&translation.AttributeKey,
		&translation.Locale,
		&translation.DisplayName,
		&translation.Unit,
		&translation.Description,
		&translation.CreatedAt,
		&translation.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("attribute translation not found")
		}
		return nil, err
	}

	return &translation, nil
}

func (r *AttributeTranslationRepository) scanTranslationFromRows(rows pgx.Rows) (*domain.AttributeTranslation, error) {
	var translation domain.AttributeTranslation

	err := rows.Scan(
		&translation.ID,
		&translation.TenantID,
		&translation.AttributeKey,
		&translation.Locale,
		&translation.DisplayName,
		&translation.Unit,
		&translation.Description,
		&translation.CreatedAt,
		&translation.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &translation, nil
}
