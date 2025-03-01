package observer

import "log"

type LoggingObserver struct{}

func (m *LoggingObserver) OnMessageProcessed(messageID string, success bool, err error) {
	if success {
		log.Printf("Message %s processed successfully", messageID)
		return
	}
	log.Printf("Message processing failed. err: %v", err)
}
