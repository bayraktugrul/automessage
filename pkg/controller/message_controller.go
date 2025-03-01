package controller

import (
	"automsg/pkg/errors"
	"automsg/pkg/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type messageController struct {
}

type MessageController interface {
	MessageSend(ctx *gin.Context)
	Messages(ctx *gin.Context)
}

func NewMessageHandler() MessageController {
	return &messageController{}
}

func (m *messageController) MessageSend(ctx *gin.Context) {
	var request request.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errors.ValidationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"operation": request.Operation,
	})
}

func (m *messageController) Messages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
