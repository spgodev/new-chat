package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
)

const limitMessages = 50

type GetChatHistoryService struct {
	messages repository.MessagesRepository
	chats    repository.DirectChatsRepository
}

func NewGetChatHistoryService(messages repository.MessagesRepository, chats repository.DirectChatsRepository) *GetChatHistoryService {
	return &GetChatHistoryService{messages: messages, chats: chats}
}

func (g *GetChatHistoryService) GetHistory(ctx context.Context,
	chatID int64,
	userID int64,
) ([]domain.Message, error) {
	chat, err := g.chats.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	isParticipant := userID == chat.User1ID || userID == chat.User2ID
	if !isParticipant {
		return nil, domain.ErrAccessDenied
	}
	messages, err := g.messages.GetByChat(ctx, chat.ID, limitMessages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
