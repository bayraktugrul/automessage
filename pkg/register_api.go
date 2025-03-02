package pkg

import (
	"automsg/docs"
	"automsg/pkg/controller"
	"automsg/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Automatic Message Sending Service
// @version 1.0
// @description API for automatic message sending system
// @BasePath /
func RegisterApi(r *gin.Engine, processControlChan chan bool, messageService service.MessageService) {

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	messageController := controller.NewMessage(processControlChan, messageService)

	r.GET("/live", func(c *gin.Context) { c.JSON(http.StatusOK, "Healthy") })
	r.GET("/messages", messageController.Messages)
	r.PUT("/send", messageController.MessageSend)
}
