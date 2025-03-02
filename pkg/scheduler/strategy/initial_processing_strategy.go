package strategy

import (
	"automsg/pkg/scheduler/observer"
	"automsg/pkg/service"
	"context"
	"log"
)

type initialProcessingStrategy struct {
	messageService    service.MessageService
	processingService service.ProcessingService
}

type ProcessingStrategy interface {
	Process(ctx context.Context, batchSize int, observerChan chan observer.Event) error
}

func NewInitialProcessingStrategy(messageService service.MessageService,
	processingService service.ProcessingService) ProcessingStrategy {

	return &initialProcessingStrategy{
		messageService:    messageService,
		processingService: processingService,
	}
}

func (i *initialProcessingStrategy) Process(ctx context.Context, batchSize int, observerChan chan observer.Event) error {
	log.Println("Initial processing is started")
	for {
		messages, err := i.messageService.GetUnsentMessages(ctx, batchSize)
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			return nil
		}

		err = i.processingService.ProcessMessages(ctx, messages, observerChan)
		if err != nil {
			return err
		}
	}
}
