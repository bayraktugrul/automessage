//go:generate go run go.uber.org/mock/mockgen -destination=../../mocks/message_service_mock.go -package=mocks automsg/pkg/service MessageService
package service

import (
	"automsg/pkg/model/dto"
	"automsg/pkg/model/response"
	"automsg/pkg/persistence"
	"context"
	"database/sql"
)

type messageService struct {
	repository persistence.MessageRepository
}

type MessageService interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error)
	GetSentMessages(ctx context.Context, page, pageSize int) ([]response.MessageResponse, int, error)

	BeginTx(ctx context.Context) (*sql.Tx, error)
	LockMessageForProcessing(ctx context.Context, tx *sql.Tx, id int64) (bool, error)
	MarkMessageAsSentTx(ctx context.Context, tx *sql.Tx, id int64, messageID string) error
}

func NewMessageService(repository persistence.MessageRepository) MessageService {
	return &messageService{
		repository: repository,
	}
}

func (s *messageService) GetUnsentMessages(ctx context.Context, limit int) ([]dto.MessageProcessingDto, error) {
	return s.repository.GetUnsentProcessingMessages(ctx, limit)
}

func (s *messageService) GetSentMessages(ctx context.Context, page, pageSize int) ([]response.MessageResponse, int, error) {
	messageDtos, totalCount, err := s.repository.GetSentMessages(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.MessageResponse, 0)
	for i := range messageDtos {
		result = append(result, messageDtos[i].ToResponse())
	}

	return result, totalCount, nil
}

func (s *messageService) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s *messageService) LockMessageForProcessing(ctx context.Context, tx *sql.Tx, id int64) (bool, error) {
	return s.repository.LockMessageForProcessing(ctx, tx, id)
}

func (s *messageService) MarkMessageAsSentTx(ctx context.Context, tx *sql.Tx, id int64, messageID string) error {
	return s.repository.MarkMessageAsSentTx(ctx, tx, id, messageID)
}
