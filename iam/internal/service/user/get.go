package user

import (
	"context"

	"iam/internal/model"
)

// GetUser возвращает пользователя по UUID
func (s *UserService) GetUser(ctx context.Context, uuid string) (*model.User, error) {
	return s.repo.GetByUUID(ctx, uuid)
}