package word

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "words"

type MongoRepository struct {
	collection *mongo.Collection
}

func NewWordMongoRepo(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoRepository) SaveWord(word *Word) (*Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	word.ID = primitive.NewObjectID()
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, word)
	if err != nil {
		return nil, err
	}

	word.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Word %s saved with ID: %s\n", word.Text, word.ID)
	return word, nil
}

func (r *MongoRepository) FindById(id string) (*Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectIDFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Fail to convert to ObjectID from id: %s", id)
		return nil, err
	}

	result := r.collection.FindOne(ctx, bson.M{"_id": objectIDFromHex})
	findWord := &Word{}
	if err := result.Decode(findWord); err != nil {
		log.Printf("Word with ID: %s not found.", id)
		return nil, err
	}

	return findWord, nil
}

func (r *MongoRepository) UpdateWord(word *Word) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	target := bson.M{"_id": word.ID}
	update := bson.M{
		"$set": bson.M{
			"text":            word.Text,
			"english_meaning": word.EnglishMeaning,
			"korean_meaning":  word.KoreanMeanings,
			"description":     word.Description,
			"synonyms":        word.Synonyms,
			"examples":        word.Examples,
			"is_delivered":    word.IsDelivered,
			"created_at":      word.CreatedAt,
			"updated_at":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, target, update)
	if err != nil {
		log.Printf("Word with ID: %s failed to update", word.ID)
		return err
	}

	return nil
}
