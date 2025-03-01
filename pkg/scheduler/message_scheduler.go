package scheduler

import (
	"automsg/pkg/scheduler/observer"
	"context"
	"log"
	"time"

	"automsg/pkg/service"
)

type MessageScheduler struct {
	messageService     service.MessageService
	processingService  service.ProcessingService
	observers          []observer.MessageObserver
	interval           time.Duration
	batchSize          int
	processControlChan chan bool
	observerChan       chan observer.Event
	done               chan struct{}
}

func NewMessageScheduler(config SchedulerConfig) *MessageScheduler {

	return &MessageScheduler{
		messageService:     config.MessageService,
		processingService:  config.ProcessingService,
		interval:           config.Interval,
		batchSize:          config.BatchSize,
		observers:          config.Observers,
		processControlChan: make(chan bool),
		observerChan:       make(chan observer.Event, 100),
		done:               make(chan struct{}),
	}
}

func (s *MessageScheduler) Start() {
	go s.run()
}

func (s *MessageScheduler) Stop() {
	s.processControlChan <- false
}

func (s *MessageScheduler) run() {
	go s.notifyObservers()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	shouldProcess := true
	log.Println("Message scheduler is starting")

	if err := s.processingService.ProcessMessages(context.Background(), s.batchSize, s.observerChan); err != nil {
		log.Printf("Error during initial processing: %v", err)
	}

	for {
		select {
		case shouldRun, ok := <-s.processControlChan:
			if !ok {
				return
			}
			shouldProcess = shouldRun
			if shouldRun {
				log.Println("Message scheduler processing started")
			} else {
				log.Println("Message scheduler processing paused")
			}
		case <-ticker.C:
			if !shouldProcess {
				continue
			}

			if err := s.processingService.ProcessMessages(context.Background(), s.batchSize, s.observerChan); err != nil {
				continue
			}

		case <-s.done:
			return
		}
	}
}

func (s *MessageScheduler) notifyObservers() {
	for {
		select {
		case evt := <-s.observerChan:
			switch evt.Type {
			case observer.EventMessageProcessed:
				for _, observer := range s.observers {
					observer.OnMessageProcessed(evt.Message.MessageID, evt.Message.Success)
				}
			}
		case <-s.done:
			return
		}
	}
}
