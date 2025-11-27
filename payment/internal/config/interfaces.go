package config

// LoggerConfig интерфейс для настроек логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// GRPCConfig интерфейс для настроек gRPC сервера
type GRPCConfig interface {
	Address() string
}

