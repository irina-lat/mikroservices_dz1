package auth

import (
	"context"
	"time"

	"iam/internal/model"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
}

// SessionRepository определяет интерфейс для работы с сессиями
type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	GetByUUID(ctx context.Context, uuid string) (*model.Session, error)
	Delete(ctx context.Context, uuid string) error
	AddToUserSet(ctx context.Context, userUUID, sessionUUID string) error
}

// Service определяет интерфейс бизнес-логики для аутентификации
type Service interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (*model.WhoamiResponse, error)
}

// AuthService реализует бизнес-логику для аутентификации
type AuthService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	sessionTTL  time.Duration
}

// NewService создаёт новый экземпляр AuthService
func NewService(
	userRepo UserRepository,
	sessionRepo SessionRepository,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}