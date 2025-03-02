package service

import (
	"automsg/pkg/client"
	"automsg/pkg/config"
	"automsg/pkg/model/dto"
	"automsg/pkg/scheduler/observer"
	"context"
	"fmt"
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

func (p *processingService) ProcessMessages(ctx context.Context, messages []dto.MessageProcessingDto, observerChan chan<- observer.Event) (err error) {
	for _, msg := range messages {
		tx, err := p.messageService.BeginTx(ctx)
		if err != nil {
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("failed to begin transaction: %w", err)},
			}
			continue
		}

		locked, err := p.messageService.LockMessageForProcessing(ctx, tx, msg.Id)
		if err != nil {
			tx.Rollback()
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("failed to lock message: %w", err)},
			}
			continue
		}

		if !locked {
			tx.Rollback()
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("message is already being processed by another instance")},
			}
			continue
		}

		resp, err := p.messageClient.SendMessage(ctx, client.Request{
			To:      msg.PhoneNumber,
			Content: msg.Content,
		})
		if err != nil {
			tx.Rollback()
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("failed to send message: %w", err)},
			}
			continue
		}

		if err := p.messageService.MarkMessageAsSentTx(ctx, tx, msg.Id, resp.MessageID); err != nil {
			tx.Rollback()
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("failed to mark message as sent: %w", err)},
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			observerChan <- observer.Event{
				Type:    observer.EventMessageProcessed,
				Message: observer.Message{Success: false, Err: fmt.Errorf("failed to commit transaction: %w", err)},
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
