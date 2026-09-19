package repository

import (
	"context"
	"time"

	"New_Chat/backend/internal/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	Revoke(ctx context.Context, tokenHash string) error
	GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	Rotate(
		ctx context.Context,
		oldTokenHash string,
		userID int64,
		newTokenHash string,
		newExpiresAt time.Time,
	) error
}

type DirectChatsRepository interface {
	Create(ctx context.Context, user1ID int64, user2ID int64, user1Title string, user2Title string) (*domain.DirectChat, error)
	GetByUser(ctx context.Context, userID int64) ([]domain.DirectChat, error)
	GetByUsers(ctx context.Context, user1ID int64, user2ID int64) (*domain.DirectChat, error)
}
