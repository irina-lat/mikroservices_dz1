package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"inventory/internal/model"
)

// FindAll возвращает все детали из MongoDB
func (r *MongoRepository) FindAll(ctx context.Context) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var parts []*model.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}
	return parts, nil
}