package story

import (
	"New_Chat/backend/internal/repository"
	"context"
)

type LogoutService struct {
	token repository.RefreshTokenRepository
}

func NewLogoutService(token repository.RefreshTokenRepository) *LogoutService {
	return &LogoutService{token: token}
}

func (l *LogoutService) Logout(ctx context.Context, token string) error {
	tokenHash := hashRefreshToken(token)
	err := l.token.Revoke(ctx, tokenHash)
	if err != nil {
		return err
	}
	return nil
}
