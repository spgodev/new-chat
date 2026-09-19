package domain

import "time"

type Message struct {
	ID       int64      `json:"id"`
	ChatsID  int64      `json:"chats_id"`
	SenderID int64      `json:"sender_id"`
	Text     string     `json:"text"`
	Status   string     `json:"status"`
	SentAt   time.Time  `json:"sent_at"`
	ReadAt   *time.Time `json:"read_at"`
}
