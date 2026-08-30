package user

import (
	"context"

	"iam/internal/model"
)

func (r *PostgresRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (uuid, login, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.UUID, user.Login, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	return err
}