package scheduler

import (
	"automsg/pkg/scheduler/observer"
	"time"
)

type SchedulerConfig struct {
	Interval           time.Duration
	Observers          []observer.MessageObserver
	ProcessControlChan chan bool
	InitialBatchSize   int
	PeriodicBatchSize  int
}
