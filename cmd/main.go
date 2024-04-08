package main

import (
	"context"
	"github.com/unbeman/av-banner-task/internal/app"
	"github.com/unbeman/av-banner-task/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	//todo: setup logger
	ctx := context.Background()

	bannerApp, err := app.GetBannerApplication(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

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
	}()

	bannerApp.Run()
}
