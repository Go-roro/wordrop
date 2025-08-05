package word

type SaveWordDto struct {
	Text           string   `json:"word" validate:"required"`
	EnglishMeaning string   `json:"english_meaning,omitempty"`
	KoreanMeanings []string `json:"korean_meaning,omitempty"`
	Description    string   `json:"description,omitempty"`
	Examples       []struct {
		ExampleText string `json:"example_text,omitempty"`
		KoreanText  string `json:"korean_text,omitempty"`
	} `json:"examples,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
}

type UpdateWordDto struct {
	ID             string   `json:"id" validate:"required"`
	Text           string   `json:"word" validate:"required"`
	EnglishMeaning string   `json:"english_meaning,omitempty"`
	KoreanMeanings []string `json:"korean_meaning,omitempty"`
	Description    string   `json:"description,omitempty"`
	Examples       []struct {
		ExampleText string `json:"example_text,omitempty"`
		KoreanText  string `json:"korean_text,omitempty"`
	} `json:"examples,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
}
