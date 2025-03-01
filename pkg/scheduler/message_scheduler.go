package scheduler

import (
	"context"
	"log"
	"time"

	"automsg/pkg/service"
)

type MessageScheduler struct {
	messageService       *service.MessageService
	interval             time.Duration
	batchSize            int
	processSchedulerChan chan bool
	doneChan             chan bool
}

func NewMessageScheduler(messageService *service.MessageService,
	interval time.Duration,
	batchSize int,
	processSchedulerChan chan bool) *MessageScheduler {

	return &MessageScheduler{
		messageService:       messageService,
		interval:             interval,
		batchSize:            batchSize,
		processSchedulerChan: processSchedulerChan,
		doneChan:             make(chan bool),
	}
}

func (s *MessageScheduler) Start() {
	go s.run()
}

func (s *MessageScheduler) Stop() {
	s.doneChan <- true
}

func (s *MessageScheduler) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	isRunning := true
	log.Println("Message scheduler is starting")

	for {
		select {
		case shouldRun, ok := <-s.processSchedulerChan:
			if !ok {
				return
			}
			isRunning = shouldRun
			if shouldRun {
				log.Println("Message scheduler processing started")
			} else {
				log.Println("Message scheduler processing paused")
			}
		case <-s.doneChan:
			return
		case <-ticker.C:
			if !isRunning {
				continue
			}
			if err := s.processUnsentMessages(); err != nil {
				log.Printf("Error processing unsent messages: %v", err)
			}
		}
	}
}

func (s *MessageScheduler) processUnsentMessages() error {
	ctx := context.Background()
	messages, err := s.messageService.GetUnsentMessages(ctx, s.batchSize)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		// TODO: Send message here to webhook url. Get webhook url from config later.
		log.Printf("Processing message ID: %d, Content: %s, Recipient: %s", msg.ID, msg.Content, msg.RecipientPhone)
		if err := s.messageService.MarkMessageAsSent(ctx, msg.ID); err != nil {
			log.Printf("Error marking message %d as sent: %v", msg.ID, err)
			continue
		}
	}

	return nil
}
