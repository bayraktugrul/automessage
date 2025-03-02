package pkg

import (
	"automsg/pkg/controller"
	"automsg/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterApi(r *gin.Engine, processControlChan chan bool, messageService service.MessageService) {

	messageHandler := controller.NewMessageHandler(processControlChan, messageService)

	r.GET("/live", func(c *gin.Context) { c.JSON(http.StatusOK, "Healthy") })
	r.GET("/messages", messageHandler.Messages)
	r.PUT("/send", messageHandler.MessageSend)
}
