package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"automsg/pkg"
	"automsg/pkg/persistence"
	"automsg/pkg/scheduler"
	"automsg/pkg/service"
)

var autoMessageApi = &cobra.Command{
	Use:   "api",
	Short: "api of automatic message sending system",
	RunE:  runApi,
}

func init() {
	root.AddCommand(autoMessageApi)
	gin.SetMode(gin.ReleaseMode)
}

func runApi(_ *cobra.Command, _ []string) error {

	db, err := persistence.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	messageRepository := persistence.NewPostgresMessageRepository(db)
	messageService := service.NewMessageService(messageRepository)

	processSchedulerChan := make(chan bool)
	messageScheduler := scheduler.NewMessageScheduler(messageService, 5*time.Second, 2, processSchedulerChan)

	r := gin.New()
	r.Use(gin.Recovery())

	pkg.RegisterApi(r)

	messageScheduler.Start()

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         ":8080",
	}

	go func() {
		log.Println("listening on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("message scheduler is shutting down...")
	messageScheduler.Stop()

	log.Println("web server shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %v", err)
	}

	return nil
}
