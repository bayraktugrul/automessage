package persistence

import (
	"context"
	"database/sql"
	"time"

	"automsg/pkg/model/dto"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{
		db: db,
	}
}

func (r *PostgresMessageRepository) GetUnsentProcessingMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error) {
	query := `
		SELECT m.id, m.content,r.phone_number
		FROM messages m
		JOIN recipients r ON m.recipient_id = r.id
		WHERE m.is_sent = false
		ORDER BY m.created_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []dto.MessageProcessingDto
	for rows.Next() {
		var msg dto.MessageProcessingDto
		err := rows.Scan(
			&msg.Id,
			&msg.Content,
			&msg.PhoneNumber,
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

func (r *PostgresMessageRepository) MarkMessageAsSent(ctx context.Context, id int64, messageID string) error {
	query := `
		UPDATE messages
		SET is_sent = true, message_id= $1, sent_at = $2, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, messageID, time.Now(), id)
	return err
}
