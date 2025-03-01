package handler

import (
	"automsg/pkg/errors"
	"automsg/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type messageHandler struct {
}

type MessageHandler interface {
	MessageSend(ctx *gin.Context)
	Messages(ctx *gin.Context)
}

func NewMessageHandler() MessageHandler {
	return &messageHandler{}
}

func (m *messageHandler) MessageSend(ctx *gin.Context) {
	var request model.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errors.ValidationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"operation": request.Operation,
	})
}

func (m *messageHandler) Messages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
