package scheduler

import (
	"automsg/mocks"
	"automsg/pkg/scheduler/observer"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func Test_MessageScheduler_should_start_and_run_initial_processing(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInitialStrategy := mocks.NewMockProcessingStrategy(ctrl)
	mockPeriodicStrategy := mocks.NewMockProcessingStrategy(ctrl)

	processControlChan := make(chan bool, 1)
	config := SchedulerConfig{
		Interval:           100 * time.Millisecond,
		ProcessControlChan: processControlChan,
		InitialBatchSize:   10,
		PeriodicBatchSize:  5,
		Observers:          []observer.MessageObserver{},
	}

	scheduler := NewMessageScheduler(mockInitialStrategy, mockPeriodicStrategy, config)

	// when
	mockInitialStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(10), gomock.Any()).Return(nil)
	mockPeriodicStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(5), gomock.Any()).Return(nil)

	// then
	scheduler.Start()
	time.Sleep(150 * time.Millisecond)
	scheduler.Stop()
}

func Test_MessageScheduler_should_handle_initial_processing_error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInitialStrategy := mocks.NewMockProcessingStrategy(ctrl)
	mockPeriodicStrategy := mocks.NewMockProcessingStrategy(ctrl)

	processControlChan := make(chan bool, 1)
	config := SchedulerConfig{
		Interval:           100 * time.Millisecond,
		ProcessControlChan: processControlChan,
		InitialBatchSize:   10,
		PeriodicBatchSize:  5,
		Observers:          []observer.MessageObserver{},
	}
	scheduler := NewMessageScheduler(mockInitialStrategy, mockPeriodicStrategy, config)

	// when
	mockInitialStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(10), gomock.Any()).
		Return(errors.New("processing error")).Times(1)
	mockPeriodicStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(5), gomock.Any()).
		Return(nil).Times(1)

	// then
	scheduler.Start()
	time.Sleep(150 * time.Millisecond)
	scheduler.Stop()
}

func Test_MessageScheduler_should_run_periodic_processing(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInitialStrategy := mocks.NewMockProcessingStrategy(ctrl)
	mockPeriodicStrategy := mocks.NewMockProcessingStrategy(ctrl)

	processControlChan := make(chan bool, 1)
	config := SchedulerConfig{
		Interval:           100 * time.Millisecond,
		ProcessControlChan: processControlChan,
		InitialBatchSize:   10,
		PeriodicBatchSize:  5,
		Observers:          []observer.MessageObserver{},
	}

	scheduler := NewMessageScheduler(mockInitialStrategy, mockPeriodicStrategy, config)

	// when
	mockInitialStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(10), gomock.Any()).Return(nil)
	mockPeriodicStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(5), gomock.Any()).Return(nil).MinTimes(1)

	// then
	scheduler.Start()
	processControlChan <- true
	time.Sleep(300 * time.Millisecond)
	scheduler.Stop()
}

func Test_MessageScheduler_should_pause_processing_when_control_channel_receives_false(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInitialStrategy := mocks.NewMockProcessingStrategy(ctrl)
	mockPeriodicStrategy := mocks.NewMockProcessingStrategy(ctrl)

	processControlChan := make(chan bool, 1)
	config := SchedulerConfig{
		Interval:           100 * time.Millisecond,
		ProcessControlChan: processControlChan,
		InitialBatchSize:   10,
		PeriodicBatchSize:  5,
		Observers:          []observer.MessageObserver{},
	}
	scheduler := NewMessageScheduler(mockInitialStrategy, mockPeriodicStrategy, config)

	// when
	mockInitialStrategy.EXPECT().Process(gomock.Any(), gomock.Eq(10), gomock.Any()).Return(nil)

	// then
	scheduler.Start()
	processControlChan <- false
	time.Sleep(300 * time.Millisecond)
	scheduler.Stop()
}
