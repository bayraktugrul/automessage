package observer

import "log"

type LoggingObserver struct{}

func (m *LoggingObserver) OnMessageProcessed(messageID int64, success bool) {
	if success {
		log.Printf("Message %d processed successfully", messageID)
		return
	}
	log.Printf("Message %d processing failed.", messageID)
}
