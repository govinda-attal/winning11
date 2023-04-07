package feeds

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbRulesProvider struct {
	db DB
}

func NewDBRulesProvider(ctx context.Context, dbURI, dbName string) (RulesProvider, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}

	return &dbRulesProvider{
		db: client.Database(dbName),
	}, nil
}

func (dp *dbRulesProvider) Rules(ctx context.Context, topic string) (*ValidationRules, error) {
	var (
		filter = bson.D{{Key: "topic", Value: topic}}
		av     ArticleValidation
		col    = dp.db.Collection("articles")
	)

	if err := col.FindOne(ctx, filter).Decode(&av); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("topic %s not found", topic)
		}

		return nil, fmt.Errorf("db rules find error: %w", err)
	}
	return &av.ValidationRules, nil
}

func (dp *dbRulesProvider) Shutdown(ctx context.Context) error {
	return dp.db.Client().Disconnect(ctx)
}
