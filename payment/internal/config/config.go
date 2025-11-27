package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/bogdanovds/rocket_factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
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

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   grpcCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию
func AppConfig() *config {
	return appConfig
}

