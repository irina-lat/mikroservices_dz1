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

	"inventory/internal/config"
	"platform/pkg/closer"
)

// App представляет основное приложение
type App struct {
	di   *DI
	name string
}

// New создаёт новое приложение
func New(cfg *config.Config) (*App, error) {
	di, err := NewDI(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		di:   di,
		name: "inventory-service",
	}, nil
}

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	log := a.di.Logger

	// 1. Запускаем gRPC сервер
	addr := a.di.Config.GRPC.Address()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Добавляем в closer для graceful shutdown
	closer.Add(func(ctx context.Context) error {
		log.Info(ctx, "Stopping gRPC server...")
		a.di.GRPCServer.GracefulStop()
		return nil
	})

	// Запускаем сервер в горутине
	go func() {
		log.Info(ctx, "Starting gRPC server", zap.String("address", addr))
		if err := a.di.GRPCServer.Serve(lis); err != nil {
			log.Error(ctx, "gRPC server error", zap.Error(err))
		}
	}()

	// 2. Ожидаем сигналы для graceful shutdown
	return a.waitForShutdown(ctx)
}

// waitForShutdown ожидает сигналы завершения
func (a *App) waitForShutdown(ctx context.Context) error {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		a.di.Logger.Info(ctx, "Context cancelled, shutting down...")
	case sig := <-shutdownCh:
		a.di.Logger.Info(ctx, "Received signal, shutting down...",
			zap.String("signal", sig.String()))
	}

	// Выполняем graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.Shutdown(shutdownCtx)
}

// Shutdown завершает работу приложения
func (a *App) Shutdown(ctx context.Context) error {
	a.di.Logger.Info(ctx, "Starting graceful shutdown...")

	done := make(chan struct{})
	go func() {
		closer.CloseAll(ctx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		a.di.Logger.Error(ctx, "Shutdown timeout exceeded, forcing exit")
		return ctx.Err()
	case <-done:
		a.di.Logger.Info(ctx, "Graceful shutdown completed")
		return nil
	}
}