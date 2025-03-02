package observer

type observerEventType int

const (
	EventMessageProcessed observerEventType = iota
)

type Event struct {
	Type    observerEventType
	Message Message
}

type Message struct {
	MessageID string
	Success   bool
	Err       error
}
