package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/bogdanovds/rocket_factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger          LoggerConfig
	HTTP            HTTPConfig
	Postgres        PostgresConfig
	InventoryClient GRPCClientConfig
	PaymentClient   GRPCClientConfig
}

// Load загружает конфигурацию из .env файла
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryClientCfg, err := env.NewInventoryClientConfig()
	if err != nil {
		return err
	}

	paymentClientCfg, err := env.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:          loggerCfg,
		HTTP:            httpCfg,
		Postgres:        postgresCfg,
		InventoryClient: inventoryClientCfg,
		PaymentClient:   paymentClientCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию
func AppConfig() *config {
	return appConfig
}

