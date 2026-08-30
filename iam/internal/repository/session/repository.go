package session

import (
	"context"
	"time"

	"iam/internal/model"
	"platform/pkg/cache"
)

type Repository interface {
	Create(ctx context.Context, session *model.Session) error
	GetByUUID(ctx context.Context, uuid string) (*model.Session, error)
	Delete(ctx context.Context, uuid string) error
	AddToUserSet(ctx context.Context, userUUID, sessionUUID string) error
	GetUserSessions(ctx context.Context, userUUID string) ([]string, error)
}

type RedisRepository struct {
	client cache.RedisClient
	ttl    time.Duration
}

func NewRedisRepository(client cache.RedisClient, ttl time.Duration) *RedisRepository {
	return &RedisRepository{
		client: client,
		ttl:    ttl,
	}
}