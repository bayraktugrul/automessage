package service

import (
	"automsg/pkg/client"
	"automsg/pkg/config"
	"automsg/pkg/model/dto"
	"automsg/pkg/scheduler/observer"
	"context"
)

type processingService struct {
	messageService MessageService
	messageClient  client.Client
	rootConfig     config.RootConfig
}

func NewProcessingService(messageService MessageService,
	messageClient client.Client,
	rootConfig config.RootConfig) ProcessingService {

	return &processingService{
		messageService: messageService,
		messageClient:  messageClient,
		rootConfig:     rootConfig,
	}
}

type ProcessingService interface {
	ProcessMessages(ctx context.Context, messages []dto.MessageProcessingDto, observerChan chan<- observer.Event) error
}

func (s *processingService) ProcessMessages(ctx context.Context, messages []dto.MessageProcessingDto, observerChan chan<- observer.Event) (err error) {
	for _, msg := range messages {
		resp, err := s.messageClient.SendMessage(ctx, client.Request{
			To:      msg.PhoneNumber,
			Content: msg.Content,
		})

		if err != nil {
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: err},
			}
			continue
		}

		if err := s.messageService.MarkMessageAsSent(ctx, msg.Id, resp.MessageID); err != nil {
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: err},
			}
			continue
		}

		observerChan <- observer.Event{
			Type:    observer.EventMessageProcessed,
			Message: observer.Message{MessageID: resp.MessageID, Success: true},
		}
	}

	return nil
}
