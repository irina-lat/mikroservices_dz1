package auth

import (
	"context"
	"time"

	"iam/internal/model"
)

// Whoami возвращает информацию о текущем пользователе по сессии
func (s *AuthService) Whoami(ctx context.Context, sessionUUID string) (*model.WhoamiResponse, error) {
	// Получаем сессию
	session, err := s.sessionRepo.GetByUUID(ctx, sessionUUID)
	if err != nil {
		return nil, model.ErrSessionNotFound
	}

	// Проверяем, не истекла ли сессия
	if session.ExpiresAt.Before(time.Now()) {
		return nil, model.ErrSessionExpired
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	return &model.WhoamiResponse{
		Session: &model.SessionInfo{
			UUID:      session.UUID,
			CreatedAt: session.CreatedAt.Format(time.RFC3339),
			UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
			ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		},
		User: user,
	}, nil
}