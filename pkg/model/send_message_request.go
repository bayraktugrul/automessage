package model

type SendMessageRequest struct {
	Operation string `json:"operation" binding:"required,oneof=START STOP"`
}
