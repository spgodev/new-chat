package story

import (
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/ws"
	"context"
)

type PresenceService struct {
	chats repository.DirectChatsRepository
	users repository.UsersRepository
	hub   *ws.Hub
}

func NewPresenceService(chats repository.DirectChatsRepository, users repository.UsersRepository, hub *ws.Hub) *PresenceService {
	return &PresenceService{chats: chats, users: users, hub: hub}
}

func (s *PresenceService) UserOnline(ctx context.Context, userID int64) error {
	chats, err := s.chats.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, chat := range chats {
		var otherUserID int64

		if userID == chat.User1ID {
			otherUserID = chat.User2ID
		} else if userID == chat.User2ID {
			otherUserID = chat.User1ID
		} else {
			continue
		}
		event := ws.Event{Type: "user_online", Data: ws.UserPresenceData{UserID: userID}}
		_ = s.hub.Send(otherUserID, event)
	}

	return nil
}

func (s *PresenceService) UserOffline(ctx context.Context, userID int64) error {
	chats, err := s.chats.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.users.UpdateLastSeen(ctx, userID); err != nil {
		return err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, chat := range chats {
		var otherUserID int64

		if userID == chat.User1ID {
			otherUserID = chat.User2ID
		} else if userID == chat.User2ID {
			otherUserID = chat.User1ID
		} else {
			continue
		}
		event := ws.Event{Type: "user_offline", Data: ws.UserPresenceData{UserID: userID, LastSeenAt: user.LastSeenAt}}
		_ = s.hub.Send(otherUserID, event)

	}

	return nil
}
