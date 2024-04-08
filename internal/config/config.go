package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
)

// Default values for config.
var (
	ServerAddressDefault      = "0.0.0.0:8080"
	PostgreSqlDSNDefault      = "postgresql://postgres:1211@localhost:5432/bkeep"
	JWTPrivateKeyDefault      = ""
	MigrationDirectoryDefault = "migrations"
)

// Config describes server's configuration, including setup for its components.
type Config struct {
	ServerAddress      string `env:"SERVER_ADDRESS"`
	PostgreSqlDSN      string `env:"POSTGRES_DSN"`
	MigrationDirectory string `env:"MIGRATION_PATH"`
	JWTPrivateKey      string `env:"JWT_PRIVATE_KEY_FILE"`
}

// parseEnv gets config setup from environment variables.
func (cfg *Config) parseEnv() error {
	return env.Parse(cfg)
}

// GetConfig returns server config.
func GetConfig() (Config, error) {
	cfg := Config{
		ServerAddress:      ServerAddressDefault,
		PostgreSqlDSN:      PostgreSqlDSNDefault,
		MigrationDirectory: MigrationDirectoryDefault,
		JWTPrivateKey:      JWTPrivateKeyDefault,
	}
	if err := cfg.parseEnv(); err != nil {
		return cfg, fmt.Errorf("could not load config from env: %w", err)
	}
	return cfg, nil
}
