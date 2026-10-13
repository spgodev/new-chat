package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/ws"
	"context"
	"errors"
	"log/slog"
	"strings"
)

type SendMessageService struct {
	messages repository.MessagesRepository
	chats    repository.DirectChatsRepository
	hub      *ws.Hub
}

func NewSendMessageService(messages repository.MessagesRepository, chats repository.DirectChatsRepository, hub *ws.Hub) *SendMessageService {
	return &SendMessageService{messages: messages, chats: chats, hub: hub}
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

	slog.Info(
		"message created",
		"message_id", message.ID,
		"chat_id", chatID,
		"sender_id", senderID,
	)

	var receiverID int64
	if senderID == chat.User1ID {
		receiverID = chat.User2ID
	} else {
		receiverID = chat.User1ID
	}

	event := ws.Event{Type: "new_message",
		Data: message}

	if err := s.hub.Send(receiverID, event); err != nil {
		slog.Debug(
			"failed to deliver new_message event",
			"message_id", message.ID,
			"chat_id", chatID,
			"receiver_id", receiverID,
			"error", err,
		)
	}

	return message, nil
}
