package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

const (
	postgresPort           = "5432"
	postgresStartupTimeout = 1 * time.Minute

	postgresEnvUserKey     = "POSTGRES_USER"
	postgresEnvPasswordKey = "POSTGRES_PASSWORD" //nolint:gosec
	postgresEnvDBKey       = "POSTGRES_DB"
)

type Container struct {
	container testcontainers.Container
	db        *sql.DB
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	dsn := buildPostgresDSN(cfg)

	db, err := connectPostgresDB(ctx, dsn)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Postgres container started", zap.String("dsn", dsn))
	success = true

	return &Container{
		container: container,
		db:        db,
		cfg:       cfg,
	}, nil
}

func (c *Container) DB() *sql.DB {
	return c.db
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) DSN() string {
	return buildPostgresDSN(c.cfg)
}

func (c *Container) URL() string {
	return buildPostgresURL(c.cfg)
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.db.Close(); err != nil {
		c.cfg.Logger.Error(ctx, "failed to close postgres connection", zap.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Postgres container terminated")

	return nil
}
