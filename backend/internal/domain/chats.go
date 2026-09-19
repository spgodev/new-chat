package domain

import "time"

type DirectChat struct {
	ID         int64     `json:"id"`
	User1ID    int64     `json:"user1_id"`
	User2ID    int64     `json:"user2_id"`
	User1Title string    `json:"user1_title"`
	User2Title string    `json:"user2_title"`
	CreatedAt  time.Time `json:"created_at"`
}
