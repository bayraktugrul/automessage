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

func Test_InitialProcessingStrategy_should_process_messages_until_empty(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewInitialProcessingStrategy(mockMessageService, mockProcessingService)

	messages1 := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
		{Id: 2, Content: "Test message 2", PhoneNumber: "+905551112244"},
	}
	emptyMessages := []dto.MessageProcessingDto{}
	observerChan := make(chan observer.Event, 10)

	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(messages1, nil)
	mockProcessingService.EXPECT().ProcessMessages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(emptyMessages, nil)

	// when
	err := strategy.Process(context.Background(), 10, observerChan)

	// then
	assert.NoError(t, err)
}

func Test_InitialProcessingStrategy_should_return_error_when_get_unsent_messages_fails(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewInitialProcessingStrategy(mockMessageService, mockProcessingService)
	observerChan := make(chan observer.Event, 10)
	expectedErr := errors.New("database error")

	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(nil, expectedErr)

	// when
	err := strategy.Process(context.Background(), 10, observerChan)

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func Test_InitialProcessingStrategy_should_return_error_when_process_messages_fails(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageService := mocks.NewMockMessageService(ctrl)
	mockProcessingService := mocks.NewMockProcessingService(ctrl)

	strategy := NewInitialProcessingStrategy(mockMessageService, mockProcessingService)

	messages := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
		{Id: 2, Content: "Test message 2", PhoneNumber: "+905551112244"},
	}

	observerChan := make(chan observer.Event, 10)
	expectedErr := errors.New("processing error")

	mockMessageService.EXPECT().GetUnsentMessages(gomock.Any(), gomock.Eq(10)).Return(messages, nil)
	mockProcessingService.EXPECT().ProcessMessages(gomock.Any(), gomock.Eq(messages), gomock.Eq(observerChan)).Return(expectedErr)

	// when
	err := strategy.Process(context.Background(), 10, observerChan)

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
