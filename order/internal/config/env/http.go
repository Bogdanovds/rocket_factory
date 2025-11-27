package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host        string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port        string `env:"HTTP_PORT" envDefault:"8081"`
	ReadTimeout string `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
}

type httpConfig struct {
	raw httpEnvConfig
}

// NewHTTPConfig создаёт конфигурацию HTTP сервера из переменных окружения
func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpConfig{raw: raw}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *httpConfig) ReadTimeout() string {
	return cfg.raw.ReadTimeout
}

