package strategy

import (
	"automsg/mocks"
	"automsg/pkg/model/dto"
	"automsg/pkg/scheduler/observer"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_PeriodicProcessingStrategy_should_process_messages(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewPeriodicProcessingStrategy(mockMessageService, mockProcessingService)

	messages := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
		{Id: 2, Content: "Test message 2", PhoneNumber: "+905551112244"},
	}

	observerChan := make(chan observer.Event, 2)
	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(2)).Return(messages, nil)
	mockProcessingService.EXPECT().ProcessMessages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// when
	err := strategy.Process(context.Background(), 2, observerChan)

	// then
	assert.NoError(t, err)
}

func Test_PeriodicProcessingStrategy_should_return_error_when_get_unsent_messages_fails(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewPeriodicProcessingStrategy(mockMessageService, mockProcessingService)

	observerChan := make(chan observer.Event, 10)
	expectedErr := errors.New("database error")

	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(nil, expectedErr)

	// when
	err := strategy.Process(context.Background(), 10, observerChan)

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func Test_PeriodicProcessingStrategy_should_return_nil_when_not_enough_messages(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewPeriodicProcessingStrategy(mockMessageService, mockProcessingService)

	messages := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
	}
	observerChan := make(chan observer.Event, 10)
	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(messages, nil)

	// when
	err := strategy.Process(context.Background(), 10, observerChan)

	// then
	assert.NoError(t, err)
}

func Test_PeriodicProcessingStrategy_should_return_error_when_process_messages_fails(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewPeriodicProcessingStrategy(mockMessageService, mockProcessingService)

	messages := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
		{Id: 2, Content: "Test message 2", PhoneNumber: "+905551112244"},
		{Id: 3, Content: "Test message 3", PhoneNumber: "+905551112255"},
		{Id: 4, Content: "Test message 4", PhoneNumber: "+905551112266"},
		{Id: 5, Content: "Test message 5", PhoneNumber: "+905551112277"},
	}
	observerChan := make(chan observer.Event, 10)
	expectedErr := errors.New("processing error")

	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(5)).Return(messages, nil)
	mockProcessingService.EXPECT().ProcessMessages(gomock.Any(), gomock.Eq(messages), gomock.Eq(observerChan)).Return(expectedErr)

	// when
	err := strategy.Process(context.Background(), 5, observerChan)

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
