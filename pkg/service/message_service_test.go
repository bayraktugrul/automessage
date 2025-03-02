package service

import (
	"automsg/mocks"
	"automsg/pkg/model/dto"
	"automsg/pkg/model/response"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_MessageService_should_get_unsent_messages(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	service := NewMessageService(mockRepo)

	expectedMessages := []dto.MessageProcessingDto{
		{Id: 1, Content: "Test message 1", PhoneNumber: "+905551112233"},
		{Id: 2, Content: "Test message 2", PhoneNumber: "+905551112244"},
	}

	// when
	mockRepo.EXPECT().
		GetUnsentProcessingMessages(gomock.Any(), gomock.Eq(10)).
		Return(expectedMessages, nil)

	// then
	result, err := service.GetUnsentMessages(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, result)
}

func Test_MessageService_should_get_sent_messages(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)
	service := NewMessageService(mockRepo)

	now := time.Now()
	messageDtos := []dto.MessageDto{
		{MessageId: "msg1", SentAt: &now},
		{MessageId: "msg2", SentAt: &now},
	}

	expectedResponses := []response.MessageResponse{
		{MessageId: "msg1", SentAt: &now},
		{MessageId: "msg2", SentAt: &now},
	}

	totalCount := 10

	// when
	mockRepo.EXPECT().
		GetSentMessages(gomock.Any(), gomock.Eq(1), gomock.Eq(20)).
		Return(messageDtos, totalCount, nil)

	// then
	result, count, err := service.GetSentMessages(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponses, result)
	assert.Equal(t, totalCount, count)
}
