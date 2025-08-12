package word

type SaveWordDto struct {
	Text           string    `json:"text" validate:"required"`
	EnglishMeaning string    `json:"english_meaning,omitempty"`
	KoreanMeanings []string  `json:"korean_meaning,omitempty"`
	Description    string    `json:"description,omitempty"`
	Examples       []Example `json:"examples,omitempty"`
	Synonyms       []string  `json:"synonyms,omitempty"`
}

type UpdateWordDto struct {
	ID             string    `json:"id" validate:"required"`
	Text           string    `json:"text" validate:"required"`
	EnglishMeaning string    `json:"english_meaning,omitempty"`
	KoreanMeanings []string  `json:"korean_meaning,omitempty"`
	Description    string    `json:"description,omitempty"`
	Examples       []Example `json:"examples,omitempty"`
	Synonyms       []string  `json:"synonyms,omitempty"`
}
