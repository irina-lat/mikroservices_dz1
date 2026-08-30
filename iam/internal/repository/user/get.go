package user

import (
	"context"
	"database/sql"

	"iam/internal/model"
)

func (r *PostgresRepository) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	query := `SELECT uuid, login, email, password, created_at, updated_at FROM users WHERE uuid = $1`
	return r.scanOne(ctx, query, uuid)
}

func (r *PostgresRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT uuid, login, email, password, created_at, updated_at FROM users WHERE login = $1`
	return r.scanOne(ctx, query, login)
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT uuid, login, email, password, created_at, updated_at FROM users WHERE email = $1`
	return r.scanOne(ctx, query, email)
}

func (r *PostgresRepository) scanOne(ctx context.Context, query string, arg string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}