package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-banner-task/internal/storage"
)

type redisManager struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisManager(redisURL string, expiration time.Duration) (storage.Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &redisManager{
		client:     client,
		expiration: expiration,
	}, nil

}

func (r redisManager) SetBanner(ctx context.Context, featureId, tagId int, bannerContent *string) error {
	key := fmt.Sprintf("%d-%d", featureId, tagId)
	err := r.client.Set(ctx, key, bannerContent, r.expiration).Err()
	if err != nil {
		return fmt.Errorf("can't exec redis set command: %w", err)
	}
	return nil
}

func (r redisManager) GetBanner(ctx context.Context, featureId, tagId int) (*string, error) {
	var bannerContent string
	key := fmt.Sprintf("%d-%d", featureId, tagId)
	bannerContent, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("no banner with key (%s): %w", key, storage.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("can't exec redis get command: %w", err)
	}
	return &bannerContent, nil
}

func (r redisManager) Shutdown() {
	r.client.Close()
	log.Info("redis client closed")
}
