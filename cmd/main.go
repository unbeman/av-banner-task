package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-banner-task/internal/app"
	"github.com/unbeman/av-banner-task/internal/config"
)

// @title Banner service
// @version 1.0
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @BasePath /
func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	level, err := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(level)

	ctx := context.Background()

	bannerApp, err := app.GetBannerApplication(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{}, 1)
	// func waits signal to stop program
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(
			exit,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)

		sig := <-exit
		log.Printf("Got signal '%v'\n", sig)

		bannerApp.Stop()
		done <- struct{}{}
	}()
	bannerApp.Run()
	<-done
}
