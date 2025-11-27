package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

// Inventory Client Config
type inventoryClientEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST" envDefault:"localhost"`
	Port string `env:"INVENTORY_GRPC_PORT" envDefault:"50051"`
}

type inventoryClientConfig struct {
	raw inventoryClientEnvConfig
}

// NewInventoryClientConfig создаёт конфигурацию клиента Inventory из переменных окружения
func NewInventoryClientConfig() (*inventoryClientConfig, error) {
	var raw inventoryClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryClientConfig{raw: raw}, nil
}

func (cfg *inventoryClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

// Payment Client Config
type paymentClientEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST" envDefault:"localhost"`
	Port string `env:"PAYMENT_GRPC_PORT" envDefault:"50052"`
}

type paymentClientConfig struct {
	raw paymentClientEnvConfig
}

// NewPaymentClientConfig создаёт конфигурацию клиента Payment из переменных окружения
func NewPaymentClientConfig() (*paymentClientConfig, error) {
	var raw paymentClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentClientConfig{raw: raw}, nil
}

func (cfg *paymentClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

