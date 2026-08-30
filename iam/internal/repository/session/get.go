package session

import (
	"context"
	"fmt"

	"iam/internal/model"
	"iam/internal/repository/converter"
)

func (r *RedisRepository) GetByUUID(ctx context.Context, uuid string) (*model.Session, error) {
	key := fmt.Sprintf("session:%s", uuid)

	data, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, model.ErrSessionNotFound
	}

	return converter.JSONToSession(data)
}

func (r *RedisRepository) Delete(ctx context.Context, uuid string) error {
	key := fmt.Sprintf("session:%s", uuid)
	return r.client.Del(ctx, key)
}