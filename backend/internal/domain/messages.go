package domain

import "time"

type MessageStatus string

const (
	MessageStatusSent MessageStatus = "sent"
	MessageStatusRead MessageStatus = "read"
)

type Message struct {
	ID       int64         `json:"id"`
	ChatID   int64         `json:"chat_id"`
	SenderID int64         `json:"sender_id"`
	Text     string        `json:"text"`
	Status   MessageStatus `json:"status"`
	SentAt   time.Time     `json:"sent_at"`
	ReadAt   *time.Time    `json:"read_at"`
}
