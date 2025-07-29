package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"wordrop/internal/domain/word"
)

const collectionName = "words"

type WordMongoRepository struct {
	collection *mongo.Collection
}

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

func NewWordMongoRepo(db *mongo.Database) word.Repository {
	return &WordMongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *WordMongoRepository) SaveWord(word *word.Word) (*word.Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, word)
	if err != nil {
		return nil, err
	}

	word.ID = result.InsertedID.(int)
	log.Printf("Word saved with ID: %s\n", word.ID)
	return word, nil
}
