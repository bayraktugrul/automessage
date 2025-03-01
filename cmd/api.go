package cmd

import (
	"automsg/pkg"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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
	r := gin.New()
	r.Use(gin.Recovery())

	pkg.RegisterApi(r)

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         ":8080",
	}

	go func() {
		fmt.Println("automsg listening on", server.Addr)
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

	fmt.Println("automessage api shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
