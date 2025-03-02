//go:generate go run go.uber.org/mock/mockgen -destination=../../../../mocks/message_observer_mock.go -package=mocks automsg/pkg/scheduler/observer MessageObserver
package observer

type MessageObserver interface {
	OnMessageProcessed(messageID string, success bool, err error)
}
