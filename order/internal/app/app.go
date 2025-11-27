package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bogdanovds/rocket_factory/order/internal/config"
	"github.com/bogdanovds/rocket_factory/platform/pkg/closer"
	"github.com/bogdanovds/rocket_factory/platform/pkg/logger"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

// App представляет приложение
type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

// New создаёт новое приложение
func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = newDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.GetLogger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderHandler := a.diContainer.OrderV1Handler(ctx)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1")
				next.ServeHTTP(w, r)
			})
		}).Mount("/", orderServer)
	})

	// Парсим read timeout
	readTimeout, err := time.ParseDuration(config.AppConfig().HTTP.ReadTimeout())
	if err != nil {
		readTimeout = 5 * time.Second
	}

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readTimeout,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().HTTP.Address()))

	lis, err := net.Listen("tcp", config.AppConfig().HTTP.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		return lis.Close()
	})

	err = a.httpServer.Serve(lis)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
