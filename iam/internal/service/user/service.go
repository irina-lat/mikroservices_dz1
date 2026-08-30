package user

import (
	"context"

	"iam/internal/model"
)

// Repository определяет интерфейс для работы с пользователями
type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

// Service определяет интерфейс бизнес-логики для пользователей
type Service interface {
	Register(ctx context.Context, login, email, password string) (*model.User, error)
	GetUser(ctx context.Context, uuid string) (*model.User, error)
}

// UserService реализует бизнес-логику для пользователей
type UserService struct {
	repo Repository
}

// NewService создаёт новый экземпляр UserService
func NewService(repo Repository) *UserService {
	return &UserService{repo: repo}
}