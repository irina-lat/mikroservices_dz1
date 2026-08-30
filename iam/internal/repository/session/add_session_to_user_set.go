package session

import (
	"context"
)

func (r *RedisRepository) AddToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	key := "user:sessions:" + userUUID
	return r.client.SAdd(ctx, key, sessionUUID)
}

func (r *RedisRepository) GetUserSessions(ctx context.Context, userUUID string) ([]string, error) {
	key := "user:sessions:" + userUUID
	return r.client.SMembers(ctx, key)
}