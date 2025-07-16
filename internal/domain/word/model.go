package word

import "time"

type Word struct {
	ID             int        `json:"id"`
	Text           string     `json:"word"`
	EnglishMeaning string     `json:"english_meaning"`
	KoreanMeaning  []string   `json:"korean_meaning"`
	Description    string     `json:"description"`
	WordExamples   *[]Example `json:"examples"`
	Synonyms       []string   `json:"synonyms"`
	IsDelivered    bool       `json:"is_delivered"`
	DeliveredAt    time.Time  `json:"delivered_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Example struct {
	ID          int    `json:"id"`
	WordID      int    `json:"word_id"`
	ExampleText string `json:"example_text"`
	KoreanText  string `json:"korean_text"`
}
