package subscription

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID               primitive.ObjectID `bson:"_id"`
	Username         string             `bson:"username"`
	Email            string             `bson:"email"`
	Verified         bool               `bson:"verified"`
	LastVerified     []time.Time        `bson:"last_verified,omitempty"`
	VerifiedCount    time.Time          `bson:"verified_time"`
	VerificationCode string             `bson:"verification_code"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}
