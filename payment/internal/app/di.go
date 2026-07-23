package app

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "payment/internal/api/payment/v1"
	"payment/internal/config"
	"payment/internal/service/payment"
	"platform/pkg/logger"
	pb "shared/pkg/proto/payment/v1"
)

type DI struct {
	Config     *config.Config
	Service    payment.Service
	API        *apiv1.API
	GRPCServer *grpc.Server
}

func NewDI(cfg *config.Config) (*DI, error) {
	if err := logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON()); err != nil {
		return nil, err
	}
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	di.Service = payment.New()
	log.Info(ctx, "Payment service initialized")

	di.API = apiv1.NewAPI(di.Service)

	di.GRPCServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	pb.RegisterPaymentServiceServer(di.GRPCServer, di.API)

	log.Info(ctx, "gRPC server configured", zap.String("addr", cfg.Payment.Address()))

	return di, nil
}