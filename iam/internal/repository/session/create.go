package session

import (
	"context"
	"fmt"

	"iam/internal/model"
	"iam/internal/repository/converter"
)

func (r *RedisRepository) Create(ctx context.Context, session *model.Session) error {
	key := fmt.Sprintf("session:%s", session.UUID)

	data, err := converter.SessionToJSON(session)
	if err != nil {
		return err
	}

	return r.client.SetWithTTL(ctx, key, data, r.ttl)
}