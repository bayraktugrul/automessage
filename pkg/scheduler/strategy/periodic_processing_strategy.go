package strategy

import (
	"automsg/pkg/scheduler/observer"
	"automsg/pkg/service"
	"context"
	"log"
)

type periodicProcessingStrategy struct {
	messageService    service.MessageService
	processingService service.ProcessingService
}

func NewPeriodicProcessingStrategy(messageService service.MessageService,
	processingService service.ProcessingService) ProcessingStrategy {

	return &periodicProcessingStrategy{
		messageService:    messageService,
		processingService: processingService}
}

func (i *periodicProcessingStrategy) Process(ctx context.Context, batchSize int, observerChan chan observer.Event) error {
	log.Println("Periodic processing is started")

	messages, err := i.messageService.GetUnsentMessages(ctx, batchSize)
	if err != nil {
		return err
	}
	if len(messages) != batchSize {
		return nil
	}

	return i.processingService.ProcessMessages(ctx, messages, observerChan)
}
