package app

import (
	"context"
	"fmt"
	"github.com/unbeman/av-banner-task/internal/config"
	controller "github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/storage"
	"github.com/unbeman/av-banner-task/internal/storage/pg"
	"github.com/unbeman/av-banner-task/internal/storage/redis"
	"github.com/unbeman/av-banner-task/internal/utils"
)

type BannerApplication struct {
	database storage.Database
	cache    storage.Cache
	server   *HTTPServer
}

func (s BannerApplication) Run() error {
	// todo: wait group
	s.server.Run()
	return nil
}

func (s BannerApplication) Stop() error {
	// todo: graceful shutdown
	s.server.Close()
	s.cache.Shutdown()
	s.database.Shutdown()
	return nil
}

func GetBannerApplication(ctx context.Context, cfg config.Config) (*BannerApplication, error) {
	pg, err := pg.NewPG(ctx, cfg.PostgreSqlDSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	redisManager, err := redis.NewRedisManager(cfg.RedisURl, cfg.RedisExpirationDuration)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	jwtManager, err := utils.NewJWTManager(cfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	ctrl, err := controller.NewController(pg, redisManager)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	hs, err := NewHTTPServer(cfg.ServerAddress, ctrl, jwtManager)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup application: %w", err)
	}

	service := &BannerApplication{
		database: pg,
		cache:    redisManager,
		server:   hs,
	}
	return service, nil
}
