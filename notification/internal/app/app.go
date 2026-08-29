package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"notification/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
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

	// Запускаем оба consumer
	go func() {
		log.Info(ctx, "Starting OrderPaid consumer...")
		if err := a.di.OrderPaidConsumer.Start(ctx); err != nil {
			log.Error(ctx, "OrderPaid consumer error", zap.Error(err))
		}
	}()

	go func() {
		log.Info(ctx, "Starting OrderAssembled consumer...")
		if err := a.di.OrderAssembledConsumer.Start(ctx); err != nil {
			log.Error(ctx, "OrderAssembled consumer error", zap.Error(err))
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