package word

import (
	"log"
	"testing"

	"github.com/Go-roro/wordrop/internal/domain/infra/db"
	"github.com/stretchr/testify/suite"
)

type WordRepoTestSuite struct {
	suite.Suite
	database *db.TestDatabase
	repo     *MongoRepository
}

func (suite *WordRepoTestSuite) SetupSuite() {
	log.Println("Setting up WordRepoTestSuite...")
	suite.database = db.SetupTestDatabase()
	suite.repo = NewWordMongoRepo(suite.database.DbInstance)
}

func (suite *WordRepoTestSuite) TearDownSuite() {
	log.Println("Tearing down WordRepoTestSuite...")
	suite.database.TearDown()
}

func (suite *WordRepoTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("Before test: %s - %s\n", suiteName, testName)
}

func TestWordRepoTestSuite(t *testing.T) {
	suite.Run(t, new(WordRepoTestSuite))
}

func (suite *WordRepoTestSuite) TestWordMongoRepository_SaveWord() {
	suite.Run("TestWordMongoRepository_SaveWord", func() {
		word := &Word{
			Text:           "test",
			EnglishMeaning: "a procedure intended to establish the quality, performance, or reliability of something",
			KoreanMeanings: []string{"테스트", "시험"},
			Description:    "A test is a method of assessing the quality or performance of something.",
			Examples: []struct {
				ExampleText string `bson:"example_text,omitempty"`
				KoreanText  string `bson:"korean_text,omitempty"`
			}{
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

		savedWord, err := suite.repo.SaveWord(word)

		suite.NotNil(savedWord.ID, "Expected saved word to have an ID")
		suite.NoError(err, "Expected no error when saving word")
		suite.NotNil(savedWord, "Expected saved word to not be nil")
	})
}

func (suite *WordRepoTestSuite) TestWordMongoRepository_FindById() {
	suite.Run("TestWordMongoRepository_FindById", func() {
		word := &Word{
			Text:           "test",
			EnglishMeaning: "a procedure intended to establish the quality, performance, or reliability of something",
			KoreanMeanings: []string{"테스트", "시험"},
			Description:    "A test is a method of assessing the quality or performance of something.",
			Examples: []struct {
				ExampleText string `bson:"example_text,omitempty"`
				KoreanText  string `bson:"korean_text,omitempty"`
			}{
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
		savedWord, err := suite.repo.SaveWord(word)
		suite.NoError(err, "Expected no error when saving word")

		findById, err := suite.repo.FindById(savedWord.ID.Hex())
		suite.NoError(err, "Expected no error when finding word by ID")
		suite.NotNil(findById, "Expected found word to not be nil")
	})
}
