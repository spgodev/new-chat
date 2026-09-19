package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	users  repository.UsersRepository
	tokens *TokenService
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

func NewRegisterService(users repository.UsersRepository, tokens *TokenService) *RegisterService {
	return &RegisterService{users: users,
		tokens: tokens}
}

func (s *RegisterService) Register(ctx context.Context, input RegisterInput) (TokenPair, error) {
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return TokenPair{}, domain.ErrInvalidEmailFormat
	}
	if len(input.Username) < 3 || len(input.Username) > 32 {
		return TokenPair{}, domain.ErrInvalidUsernameFormat
	}
	if len(input.Password) < 8 {
		return TokenPair{}, domain.ErrShortPassword
	}
	_, err = s.users.GetByUsername(ctx, input.Username)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return TokenPair{}, err
		}
	} else {
		return TokenPair{}, domain.ErrUsernameAlreadyExist
	}
	_, err = s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return TokenPair{}, err
		}
	} else {
		return TokenPair{}, domain.ErrEmailAlreadyUsed
	}
	b, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, err
	}
	user := domain.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(b),
	}
	dbUser, err := s.users.Create(ctx, user)
	if err != nil {
		return TokenPair{}, err
	}
	tokenPair, err := s.tokens.GenerateTokenPair(ctx, dbUser.ID)
	if err != nil {
		return TokenPair{}, err
	}
	return tokenPair, nil
}
