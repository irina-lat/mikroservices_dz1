package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

// Login выполняет вход пользователя и создаёт сессию
func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	// Ищем пользователя по логину или email
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		// пробуем найти по email
		user, err = s.userRepo.GetByEmail(ctx, login)
		if err != nil {
			return "", model.ErrInvalidLogin
		}
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", model.ErrInvalidPassword
	}

	// Создаём сессию
	now := time.Now()
	session := &model.Session{
		UUID:      uuid.New().String(),
		UserUUID:  user.UUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}

	// Добавляем сессию в множество сессий пользователя
	if err := s.sessionRepo.AddToUserSet(ctx, user.UUID, session.UUID); err != nil {
		// не критично
	}

	return session.UUID, nil
}