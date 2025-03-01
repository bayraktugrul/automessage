package persistence

import (
	"automsg/pkg/model/document"
	"context"
	"database/sql"
	"time"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{
		db: db,
	}
}

func (r *PostgresMessageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]document.Message, error) {
	query := `
		SELECT id, content, recipient_phone, is_sent, sent_at, created_at, updated_at
		FROM messages
		WHERE is_sent = false
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []document.Message
	for rows.Next() {
		var msg document.Message
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.RecipientPhone,
			&msg.IsSent,
			&msg.SentAt,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *PostgresMessageRepository) MarkMessageAsSent(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages
		SET is_sent = true, sent_at = $1, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), messageID)
	return err
}
