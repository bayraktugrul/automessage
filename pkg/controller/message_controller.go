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

func NewMessage(processControlChan chan<- bool,
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

// MessageSend godoc
// @Summary Send message operation
// @Description Start or stop the message sending process
// @Tags messages
// @Accept json
// @Produce json
// @Param operation body request.SendMessageRequest true "Operation details"
// @Success 200 {object} map[string]interface{} "{"operation": "START"} or {"operation": "STOP"}"
// @Failure 400 {object} errors.ErrorResponse "Validation error response"
// @Router /send [put]
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

// Messages godoc
// @Summary Get sent messages
// @Description Get a paginated list of sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param pageSize query int false "Page size (default: 10)" minimum(1) maximum(100)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} errors.ErrorResponse "Internal server error response"
// @Router /messages [get]
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
