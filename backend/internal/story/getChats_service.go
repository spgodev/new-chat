package story

import (
	"New_Chat/backend/internal/repository"
	"context"
	"time"
)

type ChatDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetChatsService struct {
	chats repository.DirectChatsRepository
}

func NewGetChatsService(chats repository.DirectChatsRepository) *GetChatsService {
	return &GetChatsService{chats: chats}
}

func (g *GetChatsService) GetByUser(ctx context.Context, userID int64) ([]ChatDTO, error) {

	chats, err := g.chats.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := []ChatDTO{}

	for _, chat := range chats {
		var title string

		if userID == chat.User1ID {
			title = chat.User1Title
		} else if userID == chat.User2ID {
			title = chat.User2Title
		} else {
			continue
		}

		result = append(result, ChatDTO{
			ID:        chat.ID,
			Title:     title,
			CreatedAt: chat.CreatedAt,
		})
	}
	return result, nil
}
