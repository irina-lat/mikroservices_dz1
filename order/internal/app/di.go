package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "order/internal/api/order/v1"
	inventoryv1 "order/internal/client/grpc/inventory/v1"
	paymentv1 "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/migrator"
	"order/internal/repository/order"
	orderservice "order/internal/service/order"
	"platform/pkg/logger"
	inventorypb "shared/pkg/proto/inventory/v1"
	paymentpb "shared/pkg/proto/payment/v1"
)

type DI struct {
	Config          *config.Config
	Pool            *pgxpool.Pool
	Repository      order.Repository
	Service         orderservice.Service
	API             *apiv1.API
	GRPCServer      *grpc.Server
	InventoryClient *inventoryv1.InventoryClient
	PaymentClient   *paymentv1.PaymentClient
}

func NewDI(cfg *config.Config) (*DI, error) {
	if err := logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON()); err != nil {
		return nil, fmt.Errorf("logger init: %w", err)
	}
	log := logger.Logger()

	di := &DI{Config: cfg}

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("db pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	di.Pool = pool
	log.Info(context.Background(), "Connected to PostgreSQL", zap.String("db", cfg.Postgres.Database()))

	// Для миграций используем *sql.DB напрямую
	sqlDB, err := sql.Open("pgx", cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	if err := migrator.Run(sqlDB); err != nil {
		log.Warn(context.Background(), "migration failed", zap.Error(err))
	}
	sqlDB.Close()

	inventoryConn, err := grpc.Dial(cfg.Inventory.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("inventory grpc dial: %w", err)
	}
	di.InventoryClient = inventoryv1.NewInventoryClient(inventorypb.NewInventoryServiceClient(inventoryConn))

	paymentConn, err := grpc.Dial(cfg.Payment.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("payment grpc dial: %w", err)
	}
	di.PaymentClient = paymentv1.NewPaymentClient(paymentpb.NewPaymentServiceClient(paymentConn))

	log.Info(context.Background(), "Connected to gRPC clients",
		zap.String("inventory", cfg.Inventory.Address()),
		zap.String("payment", cfg.Payment.Address()),
	)

	di.Repository = order.NewPostgresRepository(pool)
	di.Service = orderservice.NewService(di.Repository, di.InventoryClient, di.PaymentClient)
	di.API = apiv1.NewAPI(di.Service)
	di.GRPCServer = grpc.NewServer()

	return di, nil
}
