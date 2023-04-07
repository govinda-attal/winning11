package migrations

import (
	"context"

	"github.com/govinda-attal/winning11/feeds"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

func init() {
	ctx := context.Background()
	avs := []interface{}{
		feeds.ArticleValidation{
			Topic: "A",
			ValidationRules: feeds.ValidationRules{
				Name: feeds.ValueRule{Value: "a"},
				Desc: feeds.LenRule{LenMoreThan: 10, LenLessThan: 100},
			},
		},
		feeds.ArticleValidation{
			Topic: "B",
			ValidationRules: feeds.ValidationRules{
				Name: feeds.ValueRule{Value: "b"},
				Desc: feeds.LenRule{LenLessThan: 40},
			},
		},
	}
	err := migrate.Register(func(db *mgo.Database) error {
		_, err := db.Collection("articles").InsertMany(ctx, avs)
		if err != nil {
			return err
		}

		im := mgo.IndexModel{
			Keys: bson.D{{Key: "topic", Value: 1}},
		}
		if _, err := db.Collection("articles").Indexes().CreateOne(ctx, im); err != nil {
			return err
		}

		return nil
	}, func(db *mgo.Database) error {
		return nil
	})

	if err != nil {
		panic(err)
	}
}
