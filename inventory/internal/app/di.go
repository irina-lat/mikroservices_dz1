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
	partsvc "inventory/internal/service/part"

	"platform/pkg/closer"
	"platform/pkg/logger"
	pb "shared/pkg/proto/inventory/v1"
)

// DI контейнер для InventoryService
type DI struct {
	Config     *config.Config
	Logger     closer.Logger
	MongoDB    *mongo.Database
	Repository part.Repository
	Service    partsvc.Service
	API        *apiv1.API
	GRPCServer *grpc.Server
}

// NewDI инициализирует все зависимости
func NewDI(cfg *config.Config) (*DI, error) {
	di := &DI{
		Config: cfg,
	}

	// 1. Инициализируем логгер
	if err := di.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// 2. Инициализируем MongoDB
	if err := di.initMongoDB(); err != nil {
		return nil, fmt.Errorf("failed to init MongoDB: %w", err)
	}

	// 3. Инициализируем репозиторий
	di.initRepository()

	// 4. Инициализируем тестовые данные
	if err := di.initSampleData(); err != nil {
		di.Logger.Error(context.Background(), "Failed to init sample data", zap.Error(err))
	}

	// 5. Инициализируем сервис
	di.initService()

	// 6. Инициализируем API
	di.initAPI()

	// 7. Инициализируем gRPC сервер
	di.initGRPCServer()

	return di, nil
}

// initLogger инициализирует логгер
func (d *DI) initLogger() error {
	loggerCfg := d.Config.Logger
	if err := logger.Init(loggerCfg.Level(), loggerCfg.AsJSON()); err != nil {
		return err
	}
	d.Logger = logger.Logger()
	return nil
}

// initMongoDB инициализирует подключение к MongoDB
func (d *DI) initMongoDB() error {
	ctx := context.Background()
	mongoCfg := d.Config.Mongo

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.URI()))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	d.MongoDB = client.Database(mongoCfg.Database())
	d.Logger.Info(ctx, "Connected to MongoDB",
		zap.String("database", mongoCfg.Database()),
		zap.String("host", mongoCfg.Host()),
	)
	return nil
}

// initRepository инициализирует репозиторий
func (d *DI) initRepository() {
	d.Repository = part.NewMongoRepository(d.MongoDB)
}

// initSampleData инициализирует тестовые данные
func (d *DI) initSampleData() error {
	ctx := context.Background()
	
	// Проверяем, есть ли уже данные
	parts, err := d.Repository.FindAll(ctx)
	if err != nil {
		return err
	}
	
	if len(parts) > 0 {
		d.Logger.Info(ctx, "Sample data already exists, skipping initialization")
		return nil
	}

	d.Logger.Info(ctx, "Initializing sample data...")
	
	// Получаем репозиторий с методом InitSampleData
	repoWithInit, ok := d.Repository.(interface {
		InitSampleData(ctx context.Context) error
	})
	if !ok {
		return nil
	}
	
	if err := repoWithInit.InitSampleData(ctx); err != nil {
		return err
	}
	
	d.Logger.Info(ctx, "Sample data initialized successfully")
	return nil
}

// initService инициализирует сервисный слой
func (d *DI) initService() {
	d.Service = partsvc.NewService(d.Repository)
}

// initAPI инициализирует gRPC API
func (d *DI) initAPI() {
	d.API = apiv1.NewAPI(d.Service)
}

// initGRPCServer инициализирует gRPC сервер
func (d *DI) initGRPCServer() {
	d.GRPCServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	// Регистрируем сервис
	pb.RegisterInventoryServiceServer(d.GRPCServer, d.API)
}