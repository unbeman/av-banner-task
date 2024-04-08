package app

import (
	"context"
	"fmt"
	"github.com/unbeman/av-banner-task/internal/config"
	controller "github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/storage"
	"github.com/unbeman/av-banner-task/internal/storage/pg"
	"github.com/unbeman/av-banner-task/internal/utils"
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

func GetBannerApplication(ctx context.Context, cfg config.Config) (*BannerApplication, error) {
	pg, err := pg.NewPG(ctx, cfg.PostgreSqlDSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	jwtManager, err := utils.NewJWTManager(cfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup jwt manager: %w", err)
	}

	ctrl, err := controller.NewController(pg)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	hs, err := NewHTTPServer(cfg.ServerAddress, ctrl, jwtManager)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	service := &BannerApplication{
		storage: pg,
		server:  hs,
	}
	return service, nil
}
