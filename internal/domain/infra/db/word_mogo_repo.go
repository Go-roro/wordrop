package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"wordrop/internal/domain/word"
)

type WordMongoRepository struct {
	collection *mongo.Collection
}

func (r *WordMongoRepository) SaveWord(word *word.Word) (*word.Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	word.ID = primitive.NewObjectID()
	result, err := r.collection.InsertOne(ctx, word)
	if err != nil {
		return nil, err
	}

	word.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Word %s saved with ID: %s\n", word.Text, word.ID)
	return word, nil
}
