package dto

import (
	"github.com/Go-roro/wordrop/internal/domain/word"
)

type SaveWordRequest struct {
	Text           string         `json:"word" validate:"required"`
	EnglishMeaning string         `json:"english_meaning,omitempty"`
	KoreanMeanings []string       `json:"korean_meaning,omitempty"`
	Description    string         `json:"description,omitempty"`
	Examples       []word.Example `json:"examples,omitempty"`
	Synonyms       []string       `json:"synonyms,omitempty"`
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

type UpdateWordRequest struct {
	ID             string         `json:"id" validate:"required"`
	Text           string         `json:"word" validate:"required"`
	EnglishMeaning string         `json:"english_meaning,omitempty"`
	KoreanMeanings []string       `json:"korean_meaning,omitempty"`
	Description    string         `json:"description,omitempty"`
	Examples       []word.Example `json:"examples,omitempty"`
	Synonyms       []string       `json:"synonyms,omitempty"`
}

func (req *UpdateWordRequest) ToUpdateDto() *word.UpdateWordDto {
	return &word.UpdateWordDto{
		ID:             req.ID,
		Text:           req.Text,
		EnglishMeaning: req.EnglishMeaning,
		KoreanMeanings: req.KoreanMeanings,
		Description:    req.Description,
		Examples:       req.Examples,
		Synonyms:       req.Synonyms,
	}
}
