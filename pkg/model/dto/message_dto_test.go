package dto_test

import (
	"automsg/pkg/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_MessageDto_should_map_message_dto_to_message_response(t *testing.T) {
	//given
	timeDto := time.Time{}
	messageDto := dto.MessageDto{
		MessageId: "messageId",
		SentAt:    &timeDto,
	}

	//when
	result := messageDto.ToResponse()

	//then
	assert.Equal(t, result.MessageId, messageDto.MessageId)
	assert.Equal(t, result.SentAt, messageDto.SentAt)
}
