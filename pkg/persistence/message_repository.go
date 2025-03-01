package persistence

import (
	"automsg/pkg/model/dto"
	"context"
)

type MessageRepository interface {
	GetUnsentProcessingMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error)
	MarkMessageAsSent(ctx context.Context, id int64, messageID string) error
}
