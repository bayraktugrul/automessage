//go:generate go run go.uber.org/mock/mockgen -destination=../../../../mocks/processing_strategy_mock.go -package=mocks automsg/pkg/scheduler/strategy ProcessingStrategy
package strategy

import (
	"automsg/pkg/scheduler/observer"
	"context"
)

type ProcessingStrategy interface {
	Process(ctx context.Context, batchSize int, observerChan chan observer.Event) error
}
