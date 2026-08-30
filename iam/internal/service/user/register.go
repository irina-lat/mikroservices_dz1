package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

// Register регистрирует нового пользователя
func (s *UserService) Register(ctx context.Context, login, email, password string) (*model.User, error) {
	// Проверяем, существует ли пользователь с таким логином
	existing, _ := s.repo.GetByLogin(ctx, login)
	if existing != nil {
		return nil, model.ErrUserAlreadyExists
	}

	// Проверяем, существует ли пользователь с таким email
	existing, _ = s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, model.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	user := &model.User{
		UUID:      uuid.New().String(),
		Login:     login,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}