package service

import (
	"automsg/pkg/model/document"
	"context"
	"log"

	"automsg/pkg/persistence"
)

type messageService struct {
	repository persistence.MessageRepository
}

type MessageService interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]document.Message, error)
	MarkMessageAsSent(ctx context.Context, messageID int64) error
}

func NewMessageService(repository persistence.MessageRepository) MessageService {

	return &messageService{
		repository: repository,
	}
}

func (s *messageService) GetUnsentMessages(ctx context.Context, limit int) ([]document.Message, error) {
	return s.repository.GetUnsentMessages(ctx, limit)
}

func (s *messageService) MarkMessageAsSent(ctx context.Context, messageID int64) error {
	log.Printf("message will be marked as sent, messageId: %d", messageID)
	return s.repository.MarkMessageAsSent(ctx, messageID)
}
