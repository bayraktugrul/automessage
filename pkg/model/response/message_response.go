package response

import "time"

type MessageResponse struct {
	SentAt    *time.Time `json:"sentAt,omitempty"`
	MessageId string     `json:"messageId,omitempty"`
}
