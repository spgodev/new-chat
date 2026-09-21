package repository

import (
	"New_Chat/backend/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DirectChatPG struct {
	pool *pgxpool.Pool
}

func NewDirectChatPG(pool *pgxpool.Pool) *DirectChatPG {
	return &DirectChatPG{pool: pool}
}

func (r *DirectChatPG) Create(ctx context.Context, user1ID int64, user2ID int64, user1Title string, user2Title string) (*domain.DirectChat, error) {

	var c domain.DirectChat
	var pgErr *pgconn.PgError

	err := r.pool.QueryRow(ctx, `
INSERT INTO direct_chats (user1_id, user2_id, user1_title, user2_title)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user1_id, user2_id, user1_title, user2_title, created_at`,
		user1ID, user2ID, user1Title, user2Title).Scan(&c.ID, &c.User1ID, &c.User2ID, &c.User1Title, &c.User2Title, &c.CreatedAt)
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "direct_chats_user1_id_user2_id_key" {
					return nil, domain.ErrChatAlreadyExist
				}
				return nil, domain.ErrCantCreateChat
			}
		}
		return nil, err
	}
	return &c, nil
}

func (r *DirectChatPG) GetByUser(ctx context.Context, userID int64) ([]domain.DirectChat, error) {

	rows, err := r.pool.Query(ctx, `
SELECT id, user1_id, user2_id, user1_title, user2_title, created_at
FROM direct_chats
WHERE user1_id = $1 OR user2_id = $1
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	responds := []domain.DirectChat{}
	for rows.Next() {
		var c domain.DirectChat
		if err := rows.Scan(&c.ID, &c.User1ID, &c.User2ID, &c.User1Title, &c.User2Title, &c.CreatedAt); err != nil {
			return nil, err
		}
		responds = append(responds, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return responds, nil
}

func (r *DirectChatPG) GetByUsers(ctx context.Context, user1ID int64, user2ID int64) (*domain.DirectChat, error) {
	var c domain.DirectChat

	err := r.pool.QueryRow(ctx, `
SELECT id, user1_id, user2_id, user1_title, user2_title, created_at
FROM direct_chats
WHERE user1_id = $1 AND user2_id = $2
`, user1ID, user2ID).Scan(&c.ID, &c.User1ID, &c.User2ID, &c.User1Title, &c.User2Title, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChatNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *DirectChatPG) GetByID(ctx context.Context, chatID int64) (*domain.DirectChat, error) {
	var c domain.DirectChat

	err := r.pool.QueryRow(ctx, `
SELECT id, user1_id, user2_id, user1_title, user2_title, created_at
FROM direct_chats
WHERE id = $1
`, chatID).Scan(&c.ID, &c.User1ID, &c.User2ID, &c.User1Title, &c.User2Title, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChatNotFound
		}
		return nil, err
	}
	return &c, nil
}
