package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"context"
	"errors"
	"strings"
)

type SendMessageService struct {
	messages repository.MessagesRepository
	chats    repository.DirectChatsRepository
}

func NewSendMessageService(messages repository.MessagesRepository, chats repository.DirectChatsRepository) *SendMessageService {
	return &SendMessageService{messages: messages, chats: chats}
}

func (s *SendMessageService) Send(ctx context.Context, chatID int64,
	senderID int64, text string) (*domain.Message, error) {
	chat, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, domain.ErrChatNotFound) {
			return nil, domain.ErrChatNotFound
		}
		return nil, err
	}

	isParticipant := senderID == chat.User1ID || senderID == chat.User2ID
	if !isParticipant {
		return nil, domain.ErrAccessDenied
	}
	textToSend := strings.TrimSpace(text)
	if textToSend == "" {
		return nil, domain.ErrEmptyMessage
	}
	message, err := s.messages.Create(ctx, chatID, senderID, textToSend)
	if err != nil {
		return nil, err
	}
	return message, nil
}
