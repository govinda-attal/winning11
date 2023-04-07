package cmd

import (
	"context"
	"time"

	_ "github.com/govinda-attal/winning11/migrations"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMigrations(ctx context.Context, uri string, database string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
	}()
	// Send a ping to confirm a successful connection
	db := client.Database(database)
	var result bson.M
	if err := db.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return err
	}
	migrate.SetDatabase(db)
	migrate.SetMigrationsCollection("migrations")
	return migrate.Up(migrate.AllAvailable)
}
