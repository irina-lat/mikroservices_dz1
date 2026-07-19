package order

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// Repository определяет интерфейс для работы с заказами
type Repository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByUUID(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

// PostgresRepository реализует Repository для PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository создаёт новый PostgreSQL репозиторий
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}