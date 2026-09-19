package story

import (
	"New_Chat/backend/internal/repository"
	"context"
)

type MeDTO struct {
	ID       int64
	Username string
	Email    string
}
type MeService struct {
	users repository.UsersRepository
}

func NewMeService(users repository.UsersRepository) *MeService {
	return &MeService{users: users}
}

func (m *MeService) Me(ctx context.Context, userID int64) (MeDTO, error) {
	user, err := m.users.GetByID(ctx, userID)
	if err != nil {
		return MeDTO{}, err
	}

	return MeDTO{ID: user.ID,
		Username: user.Username,
		Email:    user.Email}, nil
}

/*if user1ID == user2ID {
return nil, domain.ErrCantCreateChat
}
if user1ID > user2ID {
swap := user2ID
user2ID = user1ID
user1ID = swap
}*/
