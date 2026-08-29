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

	// 1. Запускаем Kafka Consumer
	go func() {
		log.Info(ctx, "Starting OrderAssembled consumer...")
		if err := a.di.OrderConsumer.Start(ctx); err != nil {
			log.Error(ctx, "OrderAssembled consumer error", zap.Error(err))
		}
	}()

	// 2. Создаём HTTP роутер
	router, err := orderapi.NewServer(a.di.API, orderapi.WithPathPrefix("/api/v1"))
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

	closer.Add(func(ctx context.Context) error {
		log.Info(ctx, "Stopping HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	go func() {
		log.Info(ctx, "HTTP server started", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(ctx, "HTTP server error", zap.Error(err))
		}
	}()

	return a.waitForShutdown(ctx)
}

func (a *App) waitForShutdown(ctx context.Context) error {
	log := logger.Logger()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch
	log.Info(ctx, "Received signal, shutting down...", zap.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		closer.CloseAll(shutdownCtx)
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Warn(ctx, "Shutdown timeout")
		return shutdownCtx.Err()
	case <-done:
		log.Info(ctx, "Graceful shutdown complete")
		return nil
	}
}