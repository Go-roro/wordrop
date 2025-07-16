package word

import (
	"testing"
)

func TestWord_addMeanings(t *testing.T) {
	type input struct {
		word           *Word
		englishMeaning string
		koreanMeanings *[]string
	}

	tests := []struct {
		name     string
		input    input
		expected *Word
	}{
		{
			name: "Add English and Korean meanings",
			input: input{
				word: &Word{
					ID: 1,
				},
				englishMeaning: "test",
				koreanMeanings: &[]string{"테스트", "시험"},
			},
			expected: &Word{
				ID:             1,
				Text:           "",
				EnglishMeaning: "test",
				KoreanMeanings: []string{"테스트", "시험"},
			},
		},
		{
			name: "Append Korean meanings if existing",
			input: input{
				word: &Word{
					ID:             1,
					KoreanMeanings: []string{"평가"},
				},
				englishMeaning: "test",
				koreanMeanings: &[]string{"테스트", "시험"},
			},
			expected: &Word{
				ID:             1,
				Text:           "",
				EnglishMeaning: "test",
				KoreanMeanings: []string{"평가", "테스트", "시험"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetWord := tt.input.word
			targetWord.addMeanings(tt.input.englishMeaning, tt.input.koreanMeanings)
			if targetWord.EnglishMeaning != tt.expected.EnglishMeaning {
				t.Errorf("Expected EnglishMeaning %s, got %s", tt.expected.EnglishMeaning, targetWord.EnglishMeaning)
				return
			}

			if targetWord.KoreanMeanings == nil && tt.expected.KoreanMeanings == nil {
				t.Errorf("Expected KoreanMeanings to be nil, got nil")
				return
			}

			if len(targetWord.KoreanMeanings) != len(tt.expected.KoreanMeanings) {
				t.Errorf("Expected KoreanMeanings length %d, got %d", len(tt.expected.KoreanMeanings), len(targetWord.KoreanMeanings))
				return
			}
		})
	}
}
