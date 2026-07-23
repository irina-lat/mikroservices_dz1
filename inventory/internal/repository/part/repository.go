package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"inventory/internal/model"
)

// Repository определяет интерфейс для работы с деталями
type Repository interface {
	Save(ctx context.Context, part *model.Part) error
	FindByUUID(ctx context.Context, uuid string) (*model.Part, error)
	FindAll(ctx context.Context) ([]*model.Part, error)
}

// MongoRepository реализует Repository для MongoDB
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository создаёт новый MongoDB репозиторий
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("parts"),
	}
}

// Save сохраняет деталь в MongoDB
func (r *MongoRepository) Save(ctx context.Context, part *model.Part) error {
	now := time.Now()
	if part.CreatedAt.IsZero() {
		part.CreatedAt = now
	}
	part.UpdatedAt = now

	filter := bson.M{"uuid": part.UUID}
	update := bson.M{"$set": part}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}