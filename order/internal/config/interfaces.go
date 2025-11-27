package config

// LoggerConfig интерфейс для настроек логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// HTTPConfig интерфейс для настроек HTTP сервера
type HTTPConfig interface {
	Address() string
	ReadTimeout() string
}

// PostgresConfig интерфейс для настроек PostgreSQL
type PostgresConfig interface {
	DSN() string
	Host() string
	Port() string
	User() string
	Password() string
	Database() string
	SSLMode() string
}

// GRPCClientConfig интерфейс для настроек gRPC клиента
type GRPCClientConfig interface {
	Address() string
}

