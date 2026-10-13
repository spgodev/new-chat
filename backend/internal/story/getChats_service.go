package story

import (
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/ws"
	"context"
	"time"
)

type ChatDTO struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	CreatedAt  time.Time  `json:"created_at"`
	IsOnline   bool       `json:"is_online"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

type GetChatsService struct {
	chats repository.DirectChatsRepository
	users repository.UsersRepository
	hub   *ws.Hub
}

func NewGetChatsService(chats repository.DirectChatsRepository, users repository.UsersRepository, hub *ws.Hub) *GetChatsService {
	return &GetChatsService{chats: chats, users: users, hub: hub}
}

func (g *GetChatsService) GetByUser(ctx context.Context, userID int64) ([]ChatDTO, error) {

	chats, err := g.chats.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := []ChatDTO{}

	for _, chat := range chats {
		var title string
		var otherUserID int64

		if userID == chat.User1ID {
			title = chat.User1Title
			otherUserID = chat.User2ID
		} else if userID == chat.User2ID {
			title = chat.User2Title
			otherUserID = chat.User1ID
		} else {
			continue
		}

		otherUser, err := g.users.GetByID(ctx, otherUserID)
		if err != nil {
			return nil, err
		}

		result = append(result, ChatDTO{
			ID:         chat.ID,
			Title:      title,
			CreatedAt:  chat.CreatedAt,
			IsOnline:   g.hub.IsOnline(otherUserID),
			LastSeenAt: otherUser.LastSeenAt,
		})
	}
	return result, nil
}
