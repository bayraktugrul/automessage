package pkg

import (
	"automsg/pkg/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterApi(r *gin.Engine) {
	messageHandler := controller.NewMessageHandler()

	r.GET("/live", func(c *gin.Context) { c.JSON(http.StatusOK, "Healthy") })
	r.GET("/messages", messageHandler.Messages)
	r.PUT("/send", messageHandler.MessageSend)
}
