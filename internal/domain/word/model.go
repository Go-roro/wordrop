package word

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Word struct {
	ID             primitive.ObjectID `bson:"_id"`
	Text           string             `bson:"word"`
	EnglishMeaning string             `bson:"english_meaning"`
	KoreanMeanings []string           `bson:"korean_meaning"`
	Description    string             `bson:"description"`
	WordExamples   []*Example         `bson:"examples"`
	Synonyms       []string           `bson:"synonyms"`
	IsDelivered    bool               `bson:"is_delivered"`
	DeliveredAt    time.Time          `bson:"delivered_at"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

type Example struct {
	ExampleText string `bson:"example_text"`
	KoreanText  string `bson:"korean_text"`
}
