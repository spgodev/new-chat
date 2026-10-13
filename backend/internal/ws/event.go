package ws

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type MessagesReadData struct {
	ChatID   int64 `json:"chat_id"`
	ReaderID int64 `json:"reader_id"`
}

type IncomingEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TypingData struct {
	ChatID   int64 `json:"chat_id"`
	IsTyping bool  `json:"is_typing"`
}

type UserPresenceData struct {
	UserID     int64      `json:"user_id"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

type ChatCreatedData struct {
	ChatID int64 `json:"chat_id"`
}
