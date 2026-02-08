package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

// RefreshTokenRepository implements token storage
type RefreshTokenRepository struct {
	db *DB
}

func NewRefreshTokenRepository(db *DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE id = $1
	`

	var t domain.RefreshToken
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.UserAgent, &t.IPAddress,
		&t.CreatedAt, &t.ExpiresAt, &t.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var t domain.RefreshToken
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.UserAgent, &t.IPAddress,
		&t.CreatedAt, &t.ExpiresAt, &t.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, user_agent, ip_address, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.UserAgent, token.IPAddress,
		token.CreatedAt, token.ExpiresAt,
	)
	return err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Pool.Exec(ctx, query, time.Now())
	return err
}

// PasswordResetRepository implements password reset storage
type PasswordResetRepository struct {
	db *DB
}

func NewPasswordResetRepository(db *DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	query := `
		SELECT id, user_id, token_hash, created_at, expires_at, used_at
		FROM password_resets
		WHERE token_hash = $1
	`

	var pr domain.PasswordReset
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&pr.ID, &pr.UserID, &pr.TokenHash, &pr.CreatedAt, &pr.ExpiresAt, &pr.UsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPasswordResetNotFound
		}
		return nil, err
	}

	return &pr, nil
}

func (r *PasswordResetRepository) GetLatestByUser(ctx context.Context, userID uuid.UUID) (*domain.PasswordReset, error) {
	query := `
		SELECT id, user_id, token_hash, created_at, expires_at, used_at
		FROM password_resets
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var pr domain.PasswordReset
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&pr.ID, &pr.UserID, &pr.TokenHash, &pr.CreatedAt, &pr.ExpiresAt, &pr.UsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPasswordResetNotFound
		}
		return nil, err
	}

	return &pr, nil
}

func (r *PasswordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	query := `
		INSERT INTO password_resets (id, user_id, token_hash, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		reset.ID, reset.UserID, reset.TokenHash, reset.CreatedAt, reset.ExpiresAt,
	)
	return err
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE password_resets SET used_at = NOW() WHERE id = $1 AND used_at IS NULL`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPasswordResetNotFound
	}

	return nil
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM password_resets WHERE expires_at < $1`
	_, err := r.db.Pool.Exec(ctx, query, time.Now())
	return err
}

// AuthLogRepository implements authentication logging
type AuthLogRepository struct {
	db *DB
}

func NewAuthLogRepository(db *DB) *AuthLogRepository {
	return &AuthLogRepository{db: db}
}

func (r *AuthLogRepository) Create(ctx context.Context, log *domain.AuthenticationLog) error {
	query := `
		INSERT INTO authentication_logs (id, tenant_id, user_id, event_type, ip_address, user_agent, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		log.ID, log.TenantID, log.UserID, log.EventType, log.IPAddress, log.UserAgent,
		log.Metadata, log.CreatedAt,
	)
	return err
}

func (r *AuthLogRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error) {
	query := `
		SELECT id, tenant_id, user_id, event_type, ip_address, user_agent, metadata, created_at
		FROM authentication_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.queryLogs(ctx, query, userID, limit, offset)
}

func (r *AuthLogRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error) {
	query := `
		SELECT id, tenant_id, user_id, event_type, ip_address, user_agent, metadata, created_at
		FROM authentication_logs
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	return r.queryLogs(ctx, query, tenantID, limit, offset)
}

func (r *AuthLogRepository) queryLogs(ctx context.Context, query string, args ...any) ([]domain.AuthenticationLog, error) {
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.AuthenticationLog
	for rows.Next() {
		var log domain.AuthenticationLog
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.UserID, &log.EventType,
			&log.IPAddress, &log.UserAgent, &log.Metadata, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}
