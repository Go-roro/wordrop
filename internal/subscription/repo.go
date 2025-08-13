package subscription

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "subscriptions"

type MongoRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionRepo(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoRepository) FindByEmail(email string) (*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.collection.FindOne(ctx, bson.M{"email": email})
	subscription := &Subscription{}
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrSubscriptionNotFound
	}

	if err := result.Decode(subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (r *MongoRepository) ExistsByEmail(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false
	}

	return count > 0
}

func (r *MongoRepository) SaveSubscription(subscription *Subscription) (*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	result, err := r.collection.InsertOne(ctx, subscription)
	if err != nil {
		return nil, err
	}

	subscription.ID = result.InsertedID.(primitive.ObjectID)
	return subscription, nil
}

func (r *MongoRepository) UpdateSubscription(subscription *Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscription.UpdatedAt = time.Now()
	target := bson.M{"_id": subscription.ID}
	update := bson.M{"$set": subscription}
	if _, err := r.collection.UpdateOne(ctx, target, update); err != nil {
		return fmt.Errorf("failed to execute update for subscription %s: %w", subscription.ID.Hex(), err)
	}
	return nil
}
