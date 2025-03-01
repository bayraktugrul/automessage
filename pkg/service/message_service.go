package service

import (
	"automsg/pkg/model/document"
	"context"

	"automsg/pkg/persistence"
)

type MessageService struct {
	repository persistence.MessageRepository
}

func NewMessageService(repository persistence.MessageRepository) *MessageService {
	return &MessageService{
		repository: repository,
	}
}

func (s *MessageService) GetUnsentMessages(ctx context.Context, limit int) ([]document.Message, error) {
	return s.repository.GetUnsentMessages(ctx, limit)
}

func (s *MessageService) MarkMessageAsSent(ctx context.Context, messageID int64) error {
	return s.repository.MarkMessageAsSent(ctx, messageID)
}
