package feeds

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	Collection(string, ...*options.CollectionOptions) *mongo.Collection
	Client() *mongo.Client
}

type Validator interface {
	Validate(context.Context, Article) error
}

type RulesProvider interface {
	Rules(ctx context.Context, topic string) (*ValidationRules, error)
	Shutdown(ctx context.Context) error
}
