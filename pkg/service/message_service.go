package service

import (
	"automsg/pkg/model/dto"
	"automsg/pkg/persistence"
	"context"
)

type messageService struct {
	repository persistence.MessageRepository
}

type MessageService interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error)
	MarkMessageAsSent(ctx context.Context, id int64, messageID string) error
}

func NewMessageService(repository persistence.MessageRepository) MessageService {

	return &messageService{
		repository: repository,
	}
}

func (s *messageService) GetUnsentMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error) {
	return s.repository.GetUnsentProcessingMessages(ctx, limit)
}

func (s *messageService) MarkMessageAsSent(ctx context.Context, id int64, messageID string) error {
	return s.repository.MarkMessageAsSent(ctx, id, messageID)
}
