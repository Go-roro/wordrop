package dto

import "wordrop/internal/domain/word"

type SaveWordRequest struct {
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

func (req *SaveWordRequest) ToSaveDto() *word.SaveWordDto {
	return &word.SaveWordDto{
		Text:           req.Text,
		EnglishMeaning: req.EnglishMeaning,
		KoreanMeanings: req.KoreanMeanings,
		Description:    req.Description,
		Examples:       req.Examples,
		Synonyms:       req.Synonyms,
	}
}
