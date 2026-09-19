package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
	"errors"
)

type CreateDirectChatInput struct {
	CurrentUserID int64
	OtherUsername string
}
type CreateDirectChatService struct {
	directChats repository.DirectChatsRepository
	users       repository.UsersRepository
}

func NewCreateDirectChatService(directChats repository.DirectChatsRepository, users repository.UsersRepository) *CreateDirectChatService {
	return &CreateDirectChatService{directChats: directChats, users: users}
}

func (c *CreateDirectChatService) Create(ctx context.Context, input CreateDirectChatInput) (*domain.DirectChat, error) {

	currentUser, err := c.users.GetByID(ctx, input.CurrentUserID)
	if err != nil {
		return nil, err
	}

	otherUser, err := c.users.GetByUsername(ctx, input.OtherUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	if currentUser.ID == otherUser.ID {
		return nil, domain.ErrCantCreateChat
	}
	user1ID := currentUser.ID
	user2ID := otherUser.ID

	user1Title := "chat with " + otherUser.Username
	user2Title := "chat with " + currentUser.Username

	if user1ID > user2ID {
		user1ID = otherUser.ID
		user2ID = currentUser.ID

		user1Title = "chat with " + currentUser.Username
		user2Title = "chat with " + otherUser.Username
	}
	chat, err := c.directChats.Create(ctx, user1ID, user2ID, user1Title, user2Title)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
