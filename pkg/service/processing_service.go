package service

import (
	"automsg/pkg/model/document"
	"automsg/pkg/scheduler/observer"
	"context"
)

type processingService struct {
	messageService MessageService
}

func NewProcessingService(messageService MessageService) ProcessingService {

	return &processingService{
		messageService: messageService,
	}
}

type ProcessingService interface {
	ProcessMessages(ctx context.Context, messages []document.Message, observerChan chan observer.Event) error
}

func (s *processingService) ProcessMessages(ctx context.Context, messages []document.Message, observerChan chan observer.Event) (err error) {
	for _, msg := range messages {
		// TODO: Send external api here then check response here
		if err := s.messageService.MarkMessageAsSent(ctx, msg.ID); err != nil {
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{MessageID: msg.ID, Success: false},
			}
			continue
		}
		observerChan <- observer.Event{
			Type:    observer.EventMessageProcessed,
			Message: observer.Message{MessageID: msg.ID, Success: true},
		}
	}

	return nil
}
