package user

import (
	"context"
	"database/sql" // ← добавь этот импорт

	"iam/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}