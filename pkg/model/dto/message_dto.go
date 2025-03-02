package dto

import (
	"automsg/pkg/model/response"
	"time"
)

type MessageDto struct {
	SentAt    *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	MessageId string     `json:"message_id,omitempty" db:"message_id"`
}

func (m *MessageDto) ToResponse() response.MessageResponse {
	return response.MessageResponse{MessageId: m.MessageId, SentAt: m.SentAt}
}
