package controller

import (
	"automsg/pkg/errors"
	"automsg/pkg/model/request"
	"automsg/pkg/model/response"
	"automsg/pkg/service"
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
	messageService     service.MessageService
}

type MessageController interface {
	MessageSend(ctx *gin.Context)
	Messages(ctx *gin.Context)
}

func NewMessageHandler(processControlChan chan<- bool,
	messageService service.MessageService) MessageController {

	return &messageController{
		processControlChan: processControlChan,
		processMap: map[string]bool{
			start: true,
			stop:  false,
		},
		messageService: messageService,
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
	var req request.GetMessagesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}
	messages, totalCount, err := m.messageService.GetSentMessages(ctx, req.Page, req.PageSize)
	if err != nil {
		errors.InternalServerError(ctx, err)
		return
	}

	response := response.PaginatedResponse{
		Messages:   messages,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalCount: totalCount,
	}
	ctx.JSON(http.StatusOK, response)
}
