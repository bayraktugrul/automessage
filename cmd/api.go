package cmd

import (
	"automsg/pkg/client"
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
	"automsg/pkg/config"
	"automsg/pkg/persistence"
	"automsg/pkg/scheduler"
	"automsg/pkg/scheduler/observer"
	"automsg/pkg/scheduler/strategy"
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
	rootConfig := config.LoadConfig()

	db, err := persistence.NewConnection(rootConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	messageRepository := persistence.NewPostgresMessageRepository(db)
	messageService := service.NewMessageService(messageRepository)
	messageClient := client.New(rootConfig.App.WebhookURL)
	processingService := service.NewProcessingService(messageService, messageClient, rootConfig)
	initialProcessing := strategy.NewInitialProcessingStrategy(messageService, processingService)
	periodicProcessing := strategy.NewPeriodicProcessingStrategy(messageService, processingService)
	monitor := &observer.LoggingObserver{}

	processControlChan := make(chan bool)
	schedulerConfig := scheduler.SchedulerConfig{
		Interval:           rootConfig.App.MessageConfig.IntervalSecond,
		InitialBatchSize:   rootConfig.App.MessageConfig.InitialBatchSize,
		PeriodicBatchSize:  rootConfig.App.MessageConfig.PeriodicBatchSize,
		ProcessControlChan: processControlChan,
		Observers:          []observer.MessageObserver{monitor},
	}
	messageScheduler := scheduler.NewMessageScheduler(initialProcessing, periodicProcessing, schedulerConfig)
	messageScheduler.Start()

	r := gin.New()
	r.Use(gin.Recovery())
	pkg.RegisterApi(r, processControlChan)
	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         ":" + rootConfig.App.Port,
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
