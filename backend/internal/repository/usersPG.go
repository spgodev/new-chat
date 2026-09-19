package repository

import (
	"New_Chat/backend/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepositoryPG struct {
	pool *pgxpool.Pool
}

func NewUsersRepositoryPG(pool *pgxpool.Pool) *UsersRepositoryPG {
	return &UsersRepositoryPG{pool: pool}
}

func (r *UsersRepositoryPG) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash, created_at 
FROM users
     WHERE username = $1`, username).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepositoryPG) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash, created_at 
FROM users
     WHERE email = $1`, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepositoryPG) Create(ctx context.Context, user domain.User) (domain.User, error) {
	var u domain.User
	var pgErr *pgconn.PgError
	err := r.pool.QueryRow(ctx, `
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, password_hash, created_at`, user.Username, user.Email, user.PasswordHash).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_email_key" {
					return domain.User{}, domain.ErrEmailAlreadyUsed
				}
				if pgErr.ConstraintName == "users_username_key" {
					return domain.User{}, domain.ErrUsernameAlreadyExist
				}
			}
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UsersRepositoryPG) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash, created_at 
FROM users
     WHERE id = $1`, userID).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
