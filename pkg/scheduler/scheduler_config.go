package scheduler

import (
	"automsg/pkg/scheduler/observer"
	"time"
)

type SchedulerConfig struct {
	Interval          time.Duration
	Observers         []observer.MessageObserver
	InitialBatchSize  int
	PeriodicBatchSize int
}
