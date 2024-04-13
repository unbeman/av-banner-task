package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

// Default values for config.
var (
	PostgreSqlDSNDefault           = "postgresql://postgres:password@localhost:6500/banner-keeper?sslmode=disable"
	JWTPrivateKeyDefault           = "secret-key"
	RedisURlDefault                = "redis://default:redis-password@localhost:6379/0"
	RedisExpirationDurationDefault = 5 * time.Minute
	LogLevelDefault                = "info"
)

// Config describes server's configuration, including setup for its components.
type Config struct {
	PostgreSqlDSN           string        `env:"POSTGRES_DSN"`
	JWTPrivateKey           string        `env:"JWT_PRIVATE_KEY"`
	RedisURl                string        `env:"REDIS_URL"`
	RedisExpirationDuration time.Duration `env:"REDIS_EXPIRATION_DURATION"`
	LogLevel                string        `env:"LOG_LEVEL"`
}

// parseEnv gets config setup from environment variables.
func (cfg *Config) parseEnv() error {
	return env.Parse(cfg)
}

// GetConfig returns server config.
func GetConfig() (Config, error) {
	cfg := Config{
		PostgreSqlDSN:           PostgreSqlDSNDefault,
		JWTPrivateKey:           JWTPrivateKeyDefault,
		RedisURl:                RedisURlDefault,
		RedisExpirationDuration: RedisExpirationDurationDefault,
		LogLevel:                LogLevelDefault,
	}
	if err := cfg.parseEnv(); err != nil {
		return cfg, fmt.Errorf("could not load config from env: %w", err)
	}
	return cfg, nil
}
