package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const collectionName = "words"

func NewMongoDatabase(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connected to MongoDB")
	return client.Database(dbName), nil
}

func NewWordMongoRepo(db *mongo.Database) *WordMongoRepository {
	return &WordMongoRepository{
		collection: db.Collection(collectionName),
	}
}
