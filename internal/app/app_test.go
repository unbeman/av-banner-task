package app

import (
	"context"
	"github.com/caarlos0/env/v8"
	"github.com/stretchr/testify/suite"
	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/handlers"
	"github.com/unbeman/av-banner-task/internal/storage/pg"
	"github.com/unbeman/av-banner-task/internal/storage/redis"
	"github.com/unbeman/av-banner-task/internal/utils"
	"net/http/httptest"
	"testing"
	"time"
)

type TestConfig struct {
	ServerAddress           string        `env:"TEST_SERVER_ADDRESS"`
	PostgreSqlDSN           string        `env:"TEST_POSTGRES_DSN"`
	JWTPrivateKey           string        `env:"TEST_JWT_PRIVATE_KEY"`
	RedisURl                string        `env:"TEST_REDIS_URL"`
	RedisExpirationDuration time.Duration `env:"TEST_REDIS_EXPIRATION_DURATION"`
	LogLevel                string        `env:"TEST_LOG_LEVEL"`
}

type BannerSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *BannerSuite) SetupSuite() {
	ctx := context.Background()

	cfg := TestConfig{}
	err := env.Parse(&cfg)
	s.Nil(err)

	pg, err := pg.NewPG(ctx, cfg.PostgreSqlDSN)
	s.Nil(err)

	redisManager, err := redis.NewRedisManager(cfg.RedisURl, cfg.RedisExpirationDuration)
	s.Nil(err)

	jwtManager, err := utils.NewJWTManager(cfg.JWTPrivateKey)
	s.Nil(err)

	ctrl, err := controller.NewController(pg, redisManager)
	s.Nil(err)

	h, err := handlers.NewHttpHandler(ctrl, jwtManager)

	server := httptest.NewServer(h)
	s.server = server
}

func (s *BannerSuite) TestGetUserBanner() {
	//method := "POST"
	//uri := "/user_banner"
	//contentType := "application/json"
	//
	//s.server.Client().
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

func NewJSONRequestWithAuth() {

}
