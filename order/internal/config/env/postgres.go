package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5433"`
	User     string `env:"POSTGRES_USER" envDefault:"order-service-user"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"order-service-password"`
	Database string `env:"POSTGRES_DB" envDefault:"order-service"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

// NewPostgresConfig создаёт конфигурацию PostgreSQL из переменных окружения
func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Database,
		cfg.raw.SSLMode,
	)
}

func (cfg *postgresConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *postgresConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *postgresConfig) User() string {
	return cfg.raw.User
}

func (cfg *postgresConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *postgresConfig) Database() string {
	return cfg.raw.Database
}

func (cfg *postgresConfig) SSLMode() string {
	return cfg.raw.SSLMode
}

