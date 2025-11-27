package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// List возвращает все детали из MongoDB
func (r *Repository) List(ctx context.Context) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v", cerr)
		}
	}()

	var docs []PartDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	parts := make([]*model.Part, len(docs))
	for i, doc := range docs {
		parts[i] = ToServiceModel(&doc)
	}

	return parts, nil
}
