package word

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Word struct {
	ID             primitive.ObjectID `bson:"_id"`
	Text           string             `bson:"word"`
	EnglishMeaning string             `bson:"english_meaning"`
	KoreanMeanings []string           `bson:"korean_meaning"`
	Description    string             `bson:"description"`
	Examples       []struct {
		ExampleText string `bson:"example_text,omitempty"`
		KoreanText  string `bson:"korean_text,omitempty"`
	} `bson:"examples,omitempty"`
	Synonyms    []string  `bson:"synonyms"`
	IsDelivered bool      `bson:"is_delivered"`
	DeliveredAt time.Time `bson:"delivered_at"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
