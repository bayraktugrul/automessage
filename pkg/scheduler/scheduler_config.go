package scheduler

import (
	"automsg/pkg/scheduler/observer"
	"automsg/pkg/service"
	"time"
)

type SchedulerConfig struct {
	MessageService    service.MessageService
	ProcessingService service.ProcessingService
	Interval          time.Duration
	Observers         []observer.MessageObserver
	BatchSize         int
}
