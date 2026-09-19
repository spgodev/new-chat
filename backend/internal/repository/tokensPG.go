package repository

import (
	"New_Chat/backend/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokensRepoPG struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepoPG(pool *pgxpool.Pool) *RefreshTokensRepoPG {
	return &RefreshTokensRepoPG{pool: pool}
}

func (r *RefreshTokensRepoPG) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)`, userID, tokenHash, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokensRepoPG) Revoke(ctx context.Context, tokenHash string) error {
	tag, err := r.pool.Exec(ctx, `
UPDATE refresh_tokens
SET revoked_at = now()
WHERE token_hash = $1
`, tokenHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrTokenNotFound
	}
	return nil
}

func (r *RefreshTokensRepoPG) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	t := domain.RefreshToken{}
	err := r.pool.QueryRow(ctx, `SELECT id, user_id, token_hash, created_at, expires_at, revoked_at
       FROM refresh_tokens
       WHERE token_hash = $1`, tokenHash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.CreatedAt, &t.ExpiresAt, &t.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTokenNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokensRepoPG) Rotate(
	ctx context.Context,
	oldTokenHash string,
	userID int64,
	newTokenHash string,
	newExpiresAt time.Time,
) error {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
UPDATE refresh_tokens
SET revoked_at = now()
WHERE token_hash = $1
AND revoked_at IS NULL
`, oldTokenHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrAccessDenied
	}
	_, err = tx.Exec(ctx, `
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)`, userID, newTokenHash, newExpiresAt)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
