package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService struct {
	token        repository.RefreshTokenRepository
	accessSecret string
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

type AccessClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenService(token repository.RefreshTokenRepository, accessSecret string, accessTTL time.Duration, refreshTTL time.Duration) *TokenService {
	return &TokenService{token: token,
		accessSecret: accessSecret,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
	}
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(token string) string {
	sha := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(sha[:])
	return hashToken
}

func (s *TokenService) createRefreshToken(
	ctx context.Context,
	userID int64,
) (string, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return "", domain.ErrSomethingWrong
	}
	tokenHash := hashRefreshToken(token)
	expiresAt := time.Now().Add(s.refreshTTL)
	err = s.token.Create(ctx, userID, tokenHash, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *TokenService) createAccessToken(userID int64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTTL)

	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, userID int64) (TokenPair, error) {

	accessToken, err := s.createAccessToken(userID)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.createRefreshToken(ctx, userID)
	if err != nil {
		return TokenPair{}, err
	}

	tokenPair := TokenPair{RefreshToken: refreshToken,
		AccessToken: accessToken}
	return tokenPair, nil
}

func (s *TokenService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {

	hashOldToken := hashRefreshToken(refreshToken)
	token, err := s.token.GetByHash(ctx, hashOldToken)

	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return TokenPair{}, domain.ErrAccessDenied
		}
		return TokenPair{}, err
	}
	if time.Now().After(token.ExpiresAt) || token.RevokedAt != nil {
		return TokenPair{}, domain.ErrAccessDenied
	}

	accessToken, err := s.createAccessToken(token.UserID)
	if err != nil {
		return TokenPair{}, err
	}

	newToken, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	newExpiresAt := time.Now().Add(s.refreshTTL)

	newTokenHash := hashRefreshToken(newToken)

	err = s.token.Rotate(ctx, hashOldToken, token.UserID, newTokenHash, newExpiresAt)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessToken,
		RefreshToken: newToken}, nil
}

func (s *TokenService) ParseAccessToken(token string) (int64, error) {
	claims := AccessClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, domain.ErrAccessDenied
			}
			return []byte(s.accessSecret), nil
		})
	if err != nil {
		return 0, domain.ErrAccessDenied
	}
	if !parsedToken.Valid {
		return 0, domain.ErrAccessDenied
	}
	return claims.UserID, nil
}
