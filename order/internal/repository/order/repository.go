package order

import (
	"context"
	"database/sql"

	"order/internal/model"
)

type Repository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByUUID(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}