package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"order/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
	orderapi "shared/pkg/openapi/order/v1"
)

type App struct {
	di *DI
}

func New(cfg *config.Config) (*App, error) {
	di, err := NewDI(cfg)
	if err != nil {
		return nil, err
	}
	return &App{di: di}, nil
}

func (a *App) Run(ctx context.Context) error {
	log := logger.Logger()

	router, err := orderapi.NewServer(a.di.API)
	if err != nil {
		return fmt.Errorf("router: %w", err)
	}

	addr := a.di.Config.HTTP.Address()
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  a.di.Config.HTTP.ReadTimeout(),
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	closer.Add(func(shutdownCtx context.Context) error {
		log.Info(shutdownCtx, "Stopping HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	go func() {
		log.Info(context.Background(), "HTTP server started", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(context.Background(), "HTTP server error", zap.Error(err))
		}
	}()

	return a.waitForShutdown()
}

func (a *App) waitForShutdown() error {
	log := logger.Logger()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch
	log.Info(context.Background(), "Received signal, shutting down...", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		closer.CloseAll(context.Background())
		close(done)
	}()

	select {
	case <-ctx.Done():
		log.Warn(context.Background(), "Shutdown timeout")
		return ctx.Err()
	case <-done:
		log.Info(context.Background(), "Graceful shutdown complete")
		return nil
	}
}
