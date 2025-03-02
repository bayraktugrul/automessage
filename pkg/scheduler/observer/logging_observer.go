package observer

import (
	"log"
)

type loggingObserver struct{}

func NewLoggingObserver() MessageObserver {
	return &loggingObserver{}
}

func (m *loggingObserver) OnMessageProcessed(messageID string, success bool, err error) {
	if success {
		log.Printf("Message %s processed successfully", messageID)
		return
	}
	log.Printf("Message processing failed. err: %v", err)
}
