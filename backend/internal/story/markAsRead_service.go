package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/ws"
	"context"
)

type MarkAsReadService struct {
	message repository.MessagesRepository
	chats   repository.DirectChatsRepository
	hub     *ws.Hub
}

func NewMarkAsReadService(message repository.MessagesRepository, chats repository.DirectChatsRepository, hub *ws.Hub) *MarkAsReadService {
	return &MarkAsReadService{message: message, chats: chats, hub: hub}
}

func (m *MarkAsReadService) MarkAsRead(ctx context.Context, chatID int64, userID int64,
) error {
	chat, err := m.chats.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	isParticipant := userID == chat.User1ID || userID == chat.User2ID
	if !isParticipant {
		return domain.ErrAccessDenied
	}

	err = m.message.MarkAsRead(ctx, chat.ID, userID)
	if err != nil {
		return err
	}

	var receiverID int64

	if userID == chat.User1ID {
		receiverID = chat.User2ID
	} else {
		receiverID = chat.User1ID
	}

	event := ws.Event{Type: "messages_read",
		Data: ws.MessagesReadData{
			ChatID:   chat.ID,
			ReaderID: userID,
		}}

	_ = m.hub.Send(receiverID, event)
	return nil
}
