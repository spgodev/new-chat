package repository

import (
	"New_Chat/backend/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepositoryPG struct {
	pool *pgxpool.Pool
}

func NewMessageRepositoryPG(pool *pgxpool.Pool) *MessageRepositoryPG {
	return &MessageRepositoryPG{pool: pool}
}

func (m *MessageRepositoryPG) Create(ctx context.Context, chatID int64, senderID int64, text string) (*domain.Message, error) {
	var msg domain.Message

	err := m.pool.QueryRow(ctx, `
INSERT INTO messages (chat_id, sender_id, text)
VALUES ($1, $2, $3)
RETURNING id, chat_id, sender_id, text, status, sent_at, read_at
`, chatID, senderID, text).Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Text, &msg.Status, &msg.SentAt, &msg.ReadAt)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (m *MessageRepositoryPG) GetByChat(ctx context.Context, chatID int64, limit int) ([]domain.Message, error) {

	rows, err := m.pool.Query(ctx, `
SELECT *
FROM (
SELECT id, chat_id, sender_id, text, status, sent_at, read_at
FROM messages
WHERE chat_id = $1
ORDER BY sent_at DESC
LIMIT $2
) AS recent_messages
ORDER BY sent_at ASC;
`, chatID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	messages := []domain.Message{}

	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Text, &msg.Status, &msg.SentAt, &msg.ReadAt); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageRepositoryPG) MarkAsRead(ctx context.Context, chatID int64, readerID int64) error {
	_, err := m.pool.Exec(ctx, `
UPDATE messages
SET read_at = now(), status = 'read'
WHERE chat_id = $1 AND
      sender_id<>$2 AND
      status = 'sent'
		`, chatID, readerID)
	if err != nil {
		return err
	}
	return nil
}
