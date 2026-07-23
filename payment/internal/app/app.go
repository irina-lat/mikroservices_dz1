package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"payment/internal/config"
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
	addr := a.di.Config.Payment.Address()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	closer.Add(func(ctx context.Context) error {
		log.Info(ctx, "Stopping gRPC server...")
		a.di.GRPCServer.GracefulStop()
		return nil
	})

	go func() {
		log.Info(ctx, "gRPC server started", zap.String("addr", addr))
		if err := a.di.GRPCServer.Serve(lis); err != nil {
			log.Error(ctx, "gRPC server error", zap.Error(err))
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