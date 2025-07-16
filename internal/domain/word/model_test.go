package word

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWord_addMeanings(t *testing.T) {
	type input struct {
		word           *Word
		englishMeaning string
		koreanMeanings []string
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
				koreanMeanings: []string{"테스트", "시험"},
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
				koreanMeanings: []string{"테스트", "시험"},
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

			assert.Equal(t, tt.expected.EnglishMeaning, targetWord.EnglishMeaning)
		})
	}
}

func TestWord_appendKoreanMeanings(t *testing.T) {
	tests := []struct {
		name              string
		word              *Word
		newKoreanMeanings []string
		expected          []string
	}{
		{
			name: "Append new Korean meanings to empty list",
			word: &Word{
				ID: 1,
			},
			newKoreanMeanings: []string{"테스트", "시험"},
			expected:          []string{"테스트", "시험"},
		},
		{
			name: "Append new Korean meanings to existing list",
			word: &Word{
				ID:             1,
				KoreanMeanings: []string{"평가"},
			},
			newKoreanMeanings: []string{"테스트", "시험"},
			expected:          []string{"평가", "테스트", "시험"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.word.appendKoreanMeanings(tt.newKoreanMeanings)

			assert.Equal(t, tt.expected, tt.word.KoreanMeanings)
		})
	}
}

func TestWord_addWordExamples(t *testing.T) {
	tests := []struct {
		name     string
		word     *Word
		examples []*Example
		expected []*Example
	}{
		{
			name: "Add examples to empty word",
			word: &Word{
				ID: 1,
			},
			examples: []*Example{
				{ID: 1, ExampleText: "This is a test example.", KoreanText: "이것은 테스트 예시입니다."},
				{ID: 2, ExampleText: "Another example.", KoreanText: "또 다른 예시입니다."},
			},
			expected: []*Example{
				{ID: 1, WordID: 1, ExampleText: "This is a test example.", KoreanText: "이것은 테스트 예시입니다."},
				{ID: 2, WordID: 1, ExampleText: "Another example.", KoreanText: "또 다른 예시입니다."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.word.addWordExamples(tt.examples)
			assert.Equal(t, tt.expected, tt.word.WordExamples)
		})
	}
}

func TestWord_addSynonyms(t *testing.T) {
	tests := []struct {
		name     string
		word     *Word
		synonyms []string
		expected []string
	}{
		{
			name: "Add synonyms to empty word",
			word: &Word{
				ID: 1,
			},
			synonyms: []string{"test", "exam"},
			expected: []string{"test", "exam"},
		},
		{
			name: "Append synonyms to existing list",
			word: &Word{
				ID:       1,
				Synonyms: []string{"sample"},
			},
			synonyms: []string{"test", "exam"},
			expected: []string{"sample", "test", "exam"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.word.addSynonyms(tt.synonyms)
			assert.Equal(t, tt.expected, tt.word.Synonyms)
		})
	}
}
