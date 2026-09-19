package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string
	Password string
}

type LoginService struct {
	users  repository.UsersRepository
	tokens *TokenService
}

func NewLoginService(users repository.UsersRepository, tokens *TokenService) *LoginService {
	return &LoginService{users: users,
		tokens: tokens}
}

func (l *LoginService) Login(ctx context.Context, input LoginInput) (TokenPair, error) {
	user, err := l.users.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return TokenPair{}, domain.ErrInvalidCredentials
		}
		return TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}
	tokenPair, err := l.tokens.GenerateTokenPair(ctx, user.ID)
	if err != nil {
		return TokenPair{}, err
	}
	return tokenPair, nil
}
