package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "inventory/internal/api/inventory/v1"
	"inventory/internal/config"
	"inventory/internal/repository/part"
	partservice "inventory/internal/service/part"
	"platform/pkg/logger"
	grpcAuth "platform/pkg/middleware/grpc"
	authpb "shared/pkg/proto/auth/v1"
	pb "shared/pkg/proto/inventory/v1"
)

type DI struct {
	Config     *config.Config
	MongoDB    *mongo.Database
	Repository part.Repository
	Service    partservice.Service
	API        *apiv1.API
	GRPCServer *grpc.Server
}

func NewDI(cfg *config.Config) (*DI, error) {
	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON())
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	// 1. MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI()))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	di.MongoDB = client.Database(cfg.Mongo.Database())
	log.Info(ctx, "Connected to MongoDB", zap.String("db", cfg.Mongo.Database()))

	// 2. Репозиторий
	repo := part.NewMongoRepository(di.MongoDB)
	di.Repository = repo

	// 3. Инициализация тестовых данных
	if err := repo.InitSampleData(ctx); err != nil {
		log.Warn(ctx, "failed to init sample data", zap.Error(err))
	}

	// 4. Сервис
	di.Service = partservice.NewService(di.Repository)

	// 5. API
	di.API = apiv1.NewAPI(di.Service)

	// 6. IAM клиент
	iamConn, err := grpc.Dial(cfg.IAM.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("iam grpc dial: %w", err)
	}
	iamClient := authpb.NewAuthServiceClient(iamConn)

	// 7. Auth interceptor
	authInterceptor := grpcAuth.NewAuthInterceptor(iamClient)

	// 8. gRPC сервер с interceptor и регистрацией
	di.GRPCServer = grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	pb.RegisterInventoryServiceServer(di.GRPCServer, di.API)

	log.Info(ctx, "InventoryService DI initialized")
	return di, nil
}