package observer

type MessageObserver interface {
	OnMessageProcessed(messageID string, success bool, err error)
}
