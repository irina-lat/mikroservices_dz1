package integration

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	apiv1 "inventory/internal/api/inventory/v1"
	"inventory/internal/repository/part"
	servicepart "inventory/internal/service/part"
	pb "shared/pkg/proto/inventory/v1"
)

const bufSize = 1024 * 1024

var (
	listener *bufconn.Listener
	client   pb.InventoryServiceClient
	env      *TestEnvironment
)

func SetupTest(t *testing.T) (*TestEnvironment, pb.InventoryServiceClient, func()) {
	// 1. Поднимаем тестовое окружение
	env := NewTestEnvironment(t)

	// 2. Создаём репозиторий с тестовой БД
	db := env.MongoClient.Database(env.MongoDatabase)
	repo := part.NewMongoRepository(db)

	// 3. Добавляем тестовые данные
	ctx := context.Background()
	if err := repo.InitSampleData(ctx); err != nil {
		t.Logf("warning: failed to init sample data: %v", err)
	}

	// 4. Создаём сервис и API
	svc := servicepart.NewService(repo)
	api := apiv1.NewAPI(svc)

	// 5. Создаём gRPC сервер с bufconn
	listener = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, api)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()

	// 6. Создаём клиент
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client := pb.NewInventoryServiceClient(conn)

	teardown := func() {
		conn.Close()
		grpcServer.Stop()
		listener.Close()
		env.Close()
	}

	return env, client, teardown
}