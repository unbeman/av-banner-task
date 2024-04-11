package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"time"
)

// Default values for config.
var (
	ServerAddressDefault           = "0.0.0.0:8080"
	PostgreSqlDSNDefault           = "postgresql://postgres:password@localhost:6500/banner-keeper?sslmode=disable"
	JWTPrivateKeyDefault           = "secret-key"
	RedisURlDefault                = "redis://default:redis-password@localhost:6379/0"
	RedisExpirationDurationDefault = 5 * time.Minute
)

// Config describes server's configuration, including setup for its components.
type Config struct {
	ServerAddress           string        `env:"SERVER_ADDRESS"`
	PostgreSqlDSN           string        `env:"POSTGRES_DSN"`
	JWTPrivateKey           string        `env:"JWT_PRIVATE_KEY"`
	RedisURl                string        `env:"REDIS_URL"`
	RedisExpirationDuration time.Duration `env:"REDIS_EXPIRATION_DURATION"`
}

// parseEnv gets config setup from environment variables.
func (cfg *Config) parseEnv() error {
	return env.Parse(cfg)
}

// GetConfig returns server config.
func GetConfig() (Config, error) {
	cfg := Config{
		ServerAddress:           ServerAddressDefault,
		PostgreSqlDSN:           PostgreSqlDSNDefault,
		JWTPrivateKey:           JWTPrivateKeyDefault,
		RedisURl:                RedisURlDefault,
		RedisExpirationDuration: RedisExpirationDurationDefault,
	}
	if err := cfg.parseEnv(); err != nil {
		return cfg, fmt.Errorf("could not load config from env: %w", err)
	}
	return cfg, nil
}
