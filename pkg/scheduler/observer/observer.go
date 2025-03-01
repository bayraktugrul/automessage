package observer

type MessageObserver interface {
	OnMessageProcessed(messageID int64, success bool)
}
