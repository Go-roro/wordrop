package word

import (
	"github.com/Go-roro/wordrop/internal/common"
)

type Service struct {
	repository *MongoRepository
}

func NewWordService(repo *MongoRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) SaveNewWord(saveDto *SaveWordDto) (*Word, error) {
	word := &Word{
		Text:           saveDto.Text,
		EnglishMeaning: saveDto.EnglishMeaning,
		KoreanMeanings: saveDto.KoreanMeanings,
		Description:    saveDto.Description,
		Synonyms:       saveDto.Synonyms,
		Examples: []struct {
			ExampleText string `bson:"example_text,omitempty"`
			KoreanText  string `bson:"korean_text,omitempty"`
		}(saveDto.Examples),
		IsDelivered: false,
	}

	savedWord, err := s.repository.SaveWord(word)
	if err != nil {
		return nil, err
	}

	return savedWord, nil
}

func (s *Service) UpdateWord(updateDto *UpdateWordDto) error {
	word, err := s.repository.FindById(updateDto.ID)
	if err != nil {
		return err
	}
	updateWord := &Word{
		Text:           updateDto.Text,
		EnglishMeaning: updateDto.EnglishMeaning,
		KoreanMeanings: updateDto.KoreanMeanings,
		Description:    updateDto.Description,
		Synonyms:       updateDto.Synonyms,
		Examples: []struct {
			ExampleText string `bson:"example_text,omitempty"`
			KoreanText  string `bson:"korean_text,omitempty"`
		}(updateDto.Examples),
		IsDelivered: false,
		CreatedAt:   word.CreatedAt,
	}

	err = s.repository.UpdateWord(updateWord)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FindWords(params *SearchParams) (*common.PageResult[*Word], error) {
	if params == nil {
		params = &SearchParams{}
	}
	return s.repository.FindWords(params)
}
