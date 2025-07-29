package word

import (
	"time"
)

type Service struct {
	repository Repository
}

func NewWordService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) SaveNewWord(saveDto *SaveWordDto) (*Word, error) {
	now := time.Now()
	word := &Word{
		Text:           saveDto.Text,
		EnglishMeaning: saveDto.EnglishMeaning,
		KoreanMeanings: saveDto.KoreanMeanings,
		Description:    saveDto.Description,
		Synonyms:       saveDto.Synonyms,
		IsDelivered:    false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	savedWord, err := s.repository.SaveWord(word)
	if err != nil {
		return nil, err
	}

	return savedWord, nil
}
