package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	Port string `env:"GRPC_PORT" envDefault:"50051"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

// NewGRPCConfig создаёт конфигурацию gRPC сервера из переменных окружения
func NewGRPCConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcConfig{raw: raw}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
