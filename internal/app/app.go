package app

import (
	"fmt"
	"github.com/unbeman/av-banner-task/internal/config"
	controller "github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/storage"
	"github.com/unbeman/av-banner-task/internal/storage/pg"
)

type BannerApplication struct {
	storage storage.Database
	server  *HTTPServer
}

func (s BannerApplication) Run() error {
	// todo: wait group
	s.server.Run()
	return nil
}

func (s BannerApplication) Stop() error {
	// todo: graceful shutdown
	s.server.Close()
	return nil
}

func GetBannerApplication(cfg config.Config) (*BannerApplication, error) {
	// setup db
	// setup server
	pg, err := pg.NewPG(cfg.PostgreSqlDSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}
	ctrl, err := controller.NewController(pg)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	hs, err := NewHTTPServer(cfg.ServerAddress, ctrl)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	service := &BannerApplication{
		storage: pg,
		server:  hs,
	}
	return service, nil
}
