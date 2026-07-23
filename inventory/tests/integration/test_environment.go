package integration

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestEnvironment struct {
	MongoContainer *mongodb.MongoDBContainer
	MongoClient    *mongo.Client
	MongoDatabase  string
	GrpcEndpoint   string
	ctx            context.Context
}

func NewTestEnvironment(t *testing.T) *TestEnvironment {
	ctx := context.Background()

	// 1. Запускаем MongoDB в контейнере
	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7.0.5",
	)
	if err != nil {
		t.Fatalf("failed to start mongodb container: %v", err)
	}

	// 2. Получаем строку подключения
	connectionString, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// 3. Подключаемся к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("failed to ping mongodb: %v", err)
	}

	log.Println("✅ Test MongoDB container started")

	return &TestEnvironment{
		MongoContainer: mongoContainer,
		MongoClient:    client,
		MongoDatabase:  "inventory-test",
		GrpcEndpoint:   "localhost:50051",
		ctx:            ctx,
	}
}

func (env *TestEnvironment) Close() {
	if env.MongoClient != nil {
		env.MongoClient.Disconnect(env.ctx)
	}
	if env.MongoContainer != nil {
		env.MongoContainer.Terminate(env.ctx)
	}
}