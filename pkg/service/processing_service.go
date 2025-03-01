package service

import (
	"automsg/pkg/scheduler/observer"
	"context"
	"log"
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
	ProcessMessages(ctx context.Context, batchSize int, observerChan chan observer.Event) error
}

func (s *processingService) ProcessMessages(ctx context.Context, batchSize int, observerChan chan observer.Event) error {
	log.Println("Message scheduler is processing")

	messages, err := s.messageService.GetUnsentMessages(ctx, batchSize)
	if err != nil {
		return err
	}

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
