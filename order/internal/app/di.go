package app

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "order/internal/api/order/v1"
	inventoryclient "order/internal/client/grpc/inventory/v1"
	paymentclient "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	orderrepo "order/internal/repository/order"
	orderconsumer "order/internal/service/consumer/order_consumer"
	orderservice "order/internal/service/order"
	orderproducer "order/internal/service/producer/order_producer"
	"platform/pkg/kafka/consumer"
	"platform/pkg/kafka/producer"
	"platform/pkg/logger"
	"platform/pkg/middleware/http"
	middleware "platform/pkg/middleware/kafka"
	"platform/pkg/migrator/pg"
	authpb "shared/pkg/proto/auth/v1"
	inventorypb "shared/pkg/proto/inventory/v1"
	paymentpb "shared/pkg/proto/payment/v1"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DI struct {
	Config          *config.Config
	DB              *sql.DB
	Repository      orderrepo.Repository
	Service         orderservice.Service
	API             *apiv1.API
	GRPCServer      *grpc.Server
	InventoryClient *inventoryclient.InventoryClient
	PaymentClient   *paymentclient.PaymentClient
	OrderProducer   *orderproducer.OrderProducer
	OrderConsumer   *orderconsumer.OrderAssembledConsumer
	AuthMiddleware  *http.AuthMiddleware
}

func NewDI(cfg *config.Config) (*DI, error) {
	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON())
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	// 1. PostgreSQL
	db, err := sql.Open("pgx", cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	di.DB = db
	log.Info(ctx, "Connected to PostgreSQL", zap.String("db", cfg.Postgres.Database()))

	// 2. Миграции
	m := pg.New(db, migrationsFS, "migrations")
	if err := m.Up(ctx); err != nil {
		return nil, fmt.Errorf("migrations up: %w", err)
	}
	log.Info(ctx, "Migrations applied successfully")

	// 3. IAM клиент
	iamConn, err := grpc.Dial(cfg.IAM.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("iam grpc dial: %w", err)
	}
	iamClient := authpb.NewAuthServiceClient(iamConn)

	// 4. Auth middleware
	di.AuthMiddleware = http.NewAuthMiddleware(iamClient)

	// 5. gRPC клиенты
	inventoryConn, err := grpc.Dial(cfg.Inventory.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("inventory grpc dial: %w", err)
	}
	di.InventoryClient = inventoryclient.NewInventoryClient(inventorypb.NewInventoryServiceClient(inventoryConn))

	paymentConn, err := grpc.Dial(cfg.Payment.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("payment grpc dial: %w", err)
	}
	di.PaymentClient = paymentclient.NewPaymentClient(paymentpb.NewPaymentServiceClient(paymentConn))

	log.Info(ctx, "Connected to gRPC clients",
		zap.String("inventory", cfg.Inventory.Address()),
		zap.String("payment", cfg.Payment.Address()),
	)

	// 6. Kafka Producer
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync producer: %w", err)
	}

	kafkaProducer := producer.NewProducer(syncProducer, cfg.OrderPaidProducer.Topic(), log)
	di.OrderProducer = orderproducer.NewOrderProducer(kafkaProducer, cfg.OrderPaidProducer.Topic())

	// 7. Kafka Consumer
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderAssembledConsumer.ConsumerGroup(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	kafkaConsumer := consumer.NewConsumer(
		consumerGroup,
		[]string{cfg.OrderAssembledConsumer.Topic()},
		log,
		middleware.Logging(log),
	)

	// 8. Репозиторий, Сервис, API
	di.Repository = orderrepo.NewPostgresRepository(db)
	di.Service = orderservice.NewService(
		di.Repository,
		di.InventoryClient,
		di.PaymentClient,
		di.OrderProducer,
	)
	di.API = apiv1.NewAPI(di.Service)

	// 9. Consumer
	di.OrderConsumer = orderconsumer.NewOrderAssembledConsumer(kafkaConsumer, di.Service)

	log.Info(ctx, "OrderService DI initialized",
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers()),
		zap.String("producer_topic", cfg.OrderPaidProducer.Topic()),
		zap.String("consumer_topic", cfg.OrderAssembledConsumer.Topic()),
	)

	// 10. gRPC сервер (пока не используется)
	di.GRPCServer = grpc.NewServer()

	return di, nil
}