package observer_test

import (
	"automsg/pkg/scheduler/observer"
	"testing"

	"automsg/mocks"
	"errors"
	"go.uber.org/mock/gomock"
)

func Test_CachingObserver_should_process_successfully(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	observer := observer.NewCachingObserver(mockRedis)
	messageID := "msg-123456"

	mockRedis.EXPECT().Set(gomock.Any(), "sent_message:msg-123456", gomock.Any(), gomock.Any()).Return(nil).Times(1)

	// when && then
	observer.OnMessageProcessed(messageID, true, nil)
}

func Test_CachingObserver_should_not_process_when_success_is_false(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	observer := observer.NewCachingObserver(mockRedis)
	messageID := "msg-123456"

	// when && then
	observer.OnMessageProcessed(messageID, false, errors.New("processing error"))
}

func Test_CachingObserver_should_not_process_when_messageID_is_empty(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	observer := observer.NewCachingObserver(mockRedis)

	// when && then
	observer.OnMessageProcessed("", true, nil)
}

func Test_CachingObserver_should_handle_redis_error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	observer := observer.NewCachingObserver(mockRedis)
	messageID := "msg-123456"

	mockRedis.EXPECT().
		Set(gomock.Any(), "sent_message:msg-123456", gomock.Any(), gomock.Any()).
		Return(errors.New("redis error")).
		Times(1)

	// when && then
	observer.OnMessageProcessed(messageID, true, nil)
}
