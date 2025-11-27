package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// Get получает деталь по UUID из MongoDB
func (r *Repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	filter := bson.M{"uuid": uuid}

	var doc PartDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	return ToServiceModel(&doc), nil
}
