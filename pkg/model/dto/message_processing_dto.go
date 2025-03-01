package dto

type MessageProcessingDto struct {
	Id          int64  `json:"id" db:"id"`
	Content     string `json:"content" db:"content"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}
