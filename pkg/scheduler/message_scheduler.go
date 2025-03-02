package scheduler

import (
	"automsg/pkg/scheduler/observer"
	"automsg/pkg/scheduler/strategy"
	"context"
	"log"
	"time"
)

type MessageScheduler struct {
	initialProcessing  strategy.ProcessingStrategy
	periodicProcessing strategy.ProcessingStrategy
	config             SchedulerConfig
	observerChan       chan observer.Event
	done               chan struct{}
}

func NewMessageScheduler(initialProcessing strategy.ProcessingStrategy,
	periodicProcessing strategy.ProcessingStrategy,
	config SchedulerConfig) *MessageScheduler {

	return &MessageScheduler{
		initialProcessing:  initialProcessing,
		periodicProcessing: periodicProcessing,
		config:             config,
		observerChan:       make(chan observer.Event, 100),
		done:               make(chan struct{}),
	}
}

func (s *MessageScheduler) Start() {
	go s.run()
}

func (s *MessageScheduler) Stop() {
	s.config.ProcessControlChan <- false
}

func (s *MessageScheduler) run() {
	go s.notifyObservers()

	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	if err := s.initialProcessing.Process(context.Background(), s.config.InitialBatchSize, s.observerChan); err != nil {
		log.Printf("Error during initial processing: %v", err)
	}

	shouldProcess := true
	for {
		select {
		case run, ok := <-s.config.ProcessControlChan:
			if !ok {
				return
			}
			shouldProcess = run
			if run {
				log.Println("Message scheduler processing started")
			} else {
				log.Println("Message scheduler processing paused")
			}
		case <-ticker.C:
			if !shouldProcess {
				continue
			}

			if err := s.periodicProcessing.Process(context.Background(), s.config.PeriodicBatchSize, s.observerChan); err != nil {
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
				for _, observer := range s.config.Observers {
					observer.OnMessageProcessed(evt.Message.MessageID, evt.Message.Success, evt.Message.Err)
				}
			}
		case <-s.done:
			return
		}
	}
}
