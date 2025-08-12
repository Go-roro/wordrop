package word

import (
	"log"
	"strconv"
	"testing"

	"github.com/Go-roro/wordrop/internal/infra/testhelper"
	"github.com/stretchr/testify/suite"
)

type WordRepoTestSuite struct {
	suite.Suite
	database *testhelper.TestDatabase
	repo     *MongoRepository
}

func (suite *WordRepoTestSuite) SetupSuite() {
	log.Println("Setting up WordRepoTestSuite...")
	suite.database = testhelper.SetupTestDatabase()
	suite.repo = NewWordRepo(suite.database.DbInstance)
}

func (suite *WordRepoTestSuite) TearDownSuite() {
	log.Println("Tearing down WordRepoTestSuite...")
	suite.database.TearDown()
}

func (suite *WordRepoTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("Before test: %s - %s\n", suiteName, testName)
	if err := suite.database.CleanUp(); err != nil {
		log.Fatalf("Failed to clean up database before test: %v", err)
	}
}

func TestWordRepoTestSuite(t *testing.T) {
	suite.Run(t, new(WordRepoTestSuite))
}

func wordFixture() *Word {
	return &Word{
		Text:           "test",
		EnglishMeaning: "a procedure intended to establish the quality, performance, or reliability of something",
		KoreanMeanings: []string{"테스트", "시험"},
		Description:    "A test is a method of assessing the quality or performance of something.",
		Examples: []Example{
			{
				ExampleText: "The test results were positive.",
				KoreanText:  "시험 결과는 긍정적이었습니다.",
			},
			{
				ExampleText: "She passed the driving test on her first attempt.",
				KoreanText:  "그녀는 첫 시도에서 운전 시험에 합격했습니다.",
			},
		},
		Synonyms:    []string{"exam", "assessment"},
		IsDelivered: false,
	}
}

func (suite *WordRepoTestSuite) TestWordRepository_SaveWord() {
	suite.Run("Save", func() {
		word := wordFixture()
		savedWord, err := suite.repo.SaveWord(word)

		suite.NotNil(savedWord.ID, "Expected saved word to have an ID")
		suite.NoError(err, "Expected no error when saving word")
		suite.NotNil(savedWord, "Expected saved word to not be nil")
	})
}

func (suite *WordRepoTestSuite) TestWordRepository_FindById() {
	suite.Run("Found", func() {
		word := wordFixture()
		savedWord, err := suite.repo.SaveWord(word)
		suite.NoError(err, "Expected no error when saving word")

		findById, err := suite.repo.FindById(savedWord.ID.Hex())
		suite.NoError(err, "Expected no error when finding word by ID")
		suite.NotNil(findById, "Expected found word to not be nil")
	})
}

func (suite *WordRepoTestSuite) TestWordRepository_UpdateWord() {
	suite.Run("Update", func() {
		word := wordFixture()
		savedWord, _ := suite.repo.SaveWord(word)

		savedWord.Text = "updated test"
		savedWord.EnglishMeaning = "updated meaning"
		err := suite.repo.UpdateWord(savedWord)
		suite.NoError(err, "Expected no error when updating word")

		findById, err := suite.repo.FindById(savedWord.ID.Hex())
		suite.NoError(err, "Expected no error when finding word by ID")
		suite.Equal("updated test", findById.Text, "Expected updated word text to match")
		suite.Equal("updated meaning", findById.EnglishMeaning, "Expected updated word meaning to match")
	})
}

func (suite *WordRepoTestSuite) TestFindWords_Basic() {
	suite.Run("Basic FindWords", func() {
		_, _ = suite.repo.SaveWord(wordFixture())
		_, _ = suite.repo.SaveWord(wordFixture())

		words, err := suite.repo.FindWords(&SearchParams{})
		suite.NoError(err, "Expected no error when finding words")

		suite.Equal(2, len(words.Data), "Expected to find 2 words")
		suite.Equal(1, words.Page, "Expected page number to be 1")
	})
}

func (suite *WordRepoTestSuite) TestFindWordsWithIsDeliveredFilter() {
	suite.Run("FindWords with is_delivered filter", func() {
		wordA := wordFixture()
		wordA.IsDelivered = true
		_, _ = suite.repo.SaveWord(wordA)
		_, _ = suite.repo.SaveWord(wordFixture())

		isDelivered := true
		words, err := suite.repo.FindWords(&SearchParams{IsDelivered: &isDelivered})
		suite.NoError(err, "Expected no error when finding words with is_delivered filter")

		suite.Equal(1, len(words.Data), "Expected to find 1 word with is_delivered true")
		suite.True(words.Data[0].IsDelivered)
	})
}

func (suite *WordRepoTestSuite) TestFindWordsWithPagination() {
	suite.Run("FindWords with pagination", func() {
		for i := 0; i < 5; i++ {
			w := wordFixture()
			w.Text = "test" + strconv.Itoa(i)
			_, _ = suite.repo.SaveWord(w)
		}

		page := 2
		pageSize := 2
		words, err := suite.repo.FindWords(&SearchParams{
			Page:     page,
			PageSize: pageSize,
		})
		suite.NoError(err, "Expected no error when finding words with pagination")

		suite.Equal(pageSize, len(words.Data))
		suite.Equal(page, words.Page)
		suite.Equal(3, words.LastPage)
		suite.Equal(int64(5), words.TotalSize)
	})
}
