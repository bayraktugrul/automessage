package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"automsg/pkg/model/dto"
)

type postgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &postgresMessageRepository{
		db: db,
	}
}

func (r *postgresMessageRepository) GetUnsentProcessingMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error) {
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
		err := rows.Scan(&msg.Id, &msg.Content, &msg.PhoneNumber)
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

func (r *postgresMessageRepository) GetSentMessages(ctx context.Context, page, pageSize int) ([]dto.MessageDto, int, error) {
	offset := (page - 1) * pageSize

	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM messages m
		JOIN recipients r ON m.recipient_id = r.id
		WHERE m.is_sent = true
	`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			m.sent_at, 
			m.message_id
		FROM messages m
		JOIN recipients r ON m.recipient_id = r.id
		WHERE m.is_sent = true
		ORDER BY m.sent_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []dto.MessageDto
	for rows.Next() {
		var msg dto.MessageDto
		var messageId sql.NullString
		err := rows.Scan(&msg.SentAt, &messageId)
		if err != nil {
			return nil, 0, err
		}

		if messageId.Valid {
			msg.MessageId = messageId.String
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return messages, totalCount, nil
}

func (r *postgresMessageRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (r *postgresMessageRepository) LockMessageForProcessing(ctx context.Context, tx *sql.Tx, id int64) (bool, error) {
	query := `
		SELECT id FROM messages WHERE id = $1 AND is_sent = false FOR UPDATE SKIP LOCKED
	`
	var messageID int64
	err := tx.QueryRowContext(ctx, query, id).Scan(&messageID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to lock message: %w", err)
	}
	return true, nil
}

func (r *postgresMessageRepository) MarkMessageAsSentTx(ctx context.Context, tx *sql.Tx, id int64, messageID string) error {
	query := `
		UPDATE messages
		SET is_sent = true, message_id = $1, sent_at = $2, updated_at = $2
		WHERE id = $3
	`
	_, err := tx.ExecContext(ctx, query, messageID, time.Now(), id)
	return err
}
