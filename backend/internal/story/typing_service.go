package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/ws"
	"context"
)

type TypingService struct {
	chats repository.DirectChatsRepository
	hub   *ws.Hub
}

func NewTypingService(chats repository.DirectChatsRepository, hub *ws.Hub) *TypingService {
	return &TypingService{chats: chats, hub: hub}
}
func (t *TypingService) HandleTyping(
	ctx context.Context,
	userID int64,
	chatID int64,
	isTyping bool,
) error {
	chat, err := t.chats.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	isParticipant := userID == chat.User1ID || userID == chat.User2ID
	if !isParticipant {
		return domain.ErrAccessDenied
	}

	var receiverID int64

	if userID == chat.User1ID {
		receiverID = chat.User2ID
	} else {
		receiverID = chat.User1ID
	}

	event := ws.Event{Type: "typing",
		Data: ws.TypingData{ChatID: chatID, IsTyping: isTyping}}
	_ = t.hub.Send(receiverID, event)
	return nil
}
