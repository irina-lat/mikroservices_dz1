package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"inventory/internal/model"
)

// FindByUUID находит деталь по UUID в MongoDB
func (r *MongoRepository) FindByUUID(ctx context.Context, uuid string) (*model.Part, error) {
	filter := bson.M{"uuid": uuid}

	var part model.Part
	err := r.collection.FindOne(ctx, filter).Decode(&part)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}
	return &part, nil
}