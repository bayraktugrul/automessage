package document

import "time"

type Message struct {
	ID             int64      `json:"id" db:"id"`
	Content        string     `json:"content" db:"content"`
	RecipientPhone string     `json:"recipient_phone" db:"recipient_phone"`
	IsSent         bool       `json:"is_sent" db:"is_sent"`
	SentAt         *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
