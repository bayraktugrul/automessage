package controller

import (
	"automsg/pkg/errors"
	"automsg/pkg/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	start = "START"
	stop  = "STOP"
)

type messageController struct {
	processControlChan chan<- bool
	processMap         map[string]bool
}

type MessageController interface {
	MessageSend(ctx *gin.Context)
	Messages(ctx *gin.Context)
}

func NewMessageHandler(processControlChan chan<- bool) MessageController {
	return &messageController{
		processControlChan: processControlChan,
		processMap: map[string]bool{
			start: true,
			stop:  false,
		},
	}
}

func (m *messageController) MessageSend(ctx *gin.Context) {
	var request request.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errors.ValidationError(ctx, err)
		return
	}

	m.processControlChan <- m.processMap[request.Operation]

	ctx.JSON(http.StatusOK, gin.H{
		"operation": request.Operation,
	})
}

func (m *messageController) Messages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
