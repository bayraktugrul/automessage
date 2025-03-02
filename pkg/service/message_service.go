package service

import (
	"automsg/pkg/model/dto"
	"automsg/pkg/model/response"
	"automsg/pkg/persistence"
	"context"
)

type messageService struct {
	repository persistence.MessageRepository
}

type MessageService interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error)
	MarkMessageAsSent(ctx context.Context, id int64, messageID string) error
	GetSentMessages(ctx context.Context, page, pageSize int) ([]response.MessageResponse, int, error)
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

func (s *messageService) GetSentMessages(ctx context.Context, page, pageSize int) ([]response.MessageResponse, int, error) {
	messageDtos, totalCount, err := s.repository.GetSentMessages(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.MessageResponse, 0)
	for _, messageDto := range messageDtos {
		result = append(result, messageDto.ToResponse())
	}

	return result, totalCount, nil
}
