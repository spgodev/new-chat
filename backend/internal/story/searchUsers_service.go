package story

import (
	"New_Chat/backend/internal/repository"
	"context"
	"strings"
)

const responseLimit = 20

type SearchUsersService struct {
	users repository.UsersRepository
}

type SearchUserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func NewSearchUsersService(users repository.UsersRepository) *SearchUsersService {
	return &SearchUsersService{users: users}
}

func (s *SearchUsersService) Search(ctx context.Context, prefix string) ([]SearchUserDTO, error) {

	var respond []SearchUserDTO

	normalizePrefix := strings.TrimSpace(prefix)
	if normalizePrefix == "" {
		return respond, nil
	}
	users, err := s.users.SearchByUsernamePrefix(ctx, normalizePrefix, responseLimit)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		var sr SearchUserDTO
		sr.ID = user.ID
		sr.Username = user.Username

		respond = append(respond, sr)
	}
	return respond, nil
}
