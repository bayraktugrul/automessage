package persistence

import (
	"automsg/pkg/model/document"
	"context"
)

type MessageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]document.Message, error)
	MarkMessageAsSent(ctx context.Context, messageID int64) error
}
