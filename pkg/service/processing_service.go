package service

import (
	"automsg/pkg/config"
	"automsg/pkg/model/document"
	"automsg/pkg/scheduler/observer"
	"context"
	"log"
)

type processingService struct {
	messageService MessageService
	rootConfig     config.RootConfig
}

func NewProcessingService(messageService MessageService, rootConfig config.RootConfig) ProcessingService {

	return &processingService{
		messageService: messageService,
		rootConfig:     rootConfig,
	}
}

type ProcessingService interface {
	ProcessMessages(ctx context.Context, messages []document.Message, observerChan chan observer.Event) error
}

func (s *processingService) ProcessMessages(ctx context.Context, messages []document.Message, observerChan chan observer.Event) (err error) {
	for _, msg := range messages {
		// TODO: Send external api here then check response here
		if s.rootConfig.App.WebhookURL != "" {
			log.Printf("Would send message to webhook: %s, message ID: %d", s.rootConfig.App.WebhookURL, msg.ID)
		}

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
