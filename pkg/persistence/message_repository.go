//go:generate go run go.uber.org/mock/mockgen -destination=../../mocks/message_repository_mock.go -package=mocks automsg/pkg/persistence MessageRepository
package persistence

import (
	"automsg/pkg/model/dto"
	"context"
	"database/sql"
)

type MessageRepository interface {
	GetUnsentProcessingMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error)
	GetSentMessages(ctx context.Context, page, pageSize int) ([]dto.MessageDto, int, error)

	BeginTx(ctx context.Context) (*sql.Tx, error)
	LockMessageForProcessing(ctx context.Context, tx *sql.Tx, id int64) (bool, error)
	MarkMessageAsSentTx(ctx context.Context, tx *sql.Tx, id int64, messageID string) error
}
