package subscription

import (
	"log"
	"testing"

	"github.com/Go-roro/wordrop/internal/infra/testhelper"
	"github.com/stretchr/testify/suite"
)

type SubscriptionRepoTestSuite struct {
	suite.Suite
	database *testhelper.TestDatabase
	repo     *MongoRepository
}

func (suite *SubscriptionRepoTestSuite) SetupSuite() {
	log.Println("Setting up SubscriptionRepoTestSuite...")
	suite.database = testhelper.SetupTestDatabase()
	suite.repo = NewSubscriptionRepo(suite.database.DbInstance)
}

func (suite *SubscriptionRepoTestSuite) TearDownSuite() {
	log.Println("Tearing down SubscriptionRepoTestSuite...")
	suite.database.TearDown()
}

func (suite *SubscriptionRepoTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("Before test: %s - %s\n", suiteName, testName)
	if err := suite.database.CleanUp(); err != nil {
		log.Fatalf("Failed to clean up database before test: %v", err)
	}
}

func TestSubscriptionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionRepoTestSuite))
}

func subscriptionFixture() *Subscription {
	return NewSubscription("testuser", "test@example.com")
}

func (suite *SubscriptionRepoTestSuite) TestSubscriptionRepository_SaveSubscription() {
	suite.Run("Save", func() {
		sub := subscriptionFixture()
		savedSub, err := suite.repo.SaveSubscription(sub)
		suite.NoError(err, "Expected no error when saving subscription")

		suite.NotNil(savedSub, "Expected saved subscription to not be nil")
		suite.NotNil(savedSub.ID, "Expected saved subscription to have an ID")
		suite.Equal(sub.Email, savedSub.Email)
	})
}

func (suite *SubscriptionRepoTestSuite) TestSubscriptionRepository_FindByEmail() {
	suite.Run("Found", func() {
		sub := subscriptionFixture()
		_, err := suite.repo.SaveSubscription(sub)
		suite.NoError(err)

		foundSub, err := suite.repo.FindByEmail(sub.Email)
		suite.NoError(err, "Expected no error when finding subscription by email")
		suite.NotNil(foundSub, "Expected found subscription to not be nil")
		suite.Equal(sub.Email, foundSub.Email)
	})

	suite.Run("NotFound", func() {
		_, err := suite.repo.FindByEmail("nonexistent@example.com")
		suite.Error(err, "Expected an error when subscription is not found")
	})
}

func (suite *SubscriptionRepoTestSuite) TestSubscriptionRepository_UpdateSubscription() {
	suite.Run("Update", func() {
		sub := subscriptionFixture()
		savedSub, _ := suite.repo.SaveSubscription(sub)

		savedSub.Verified = true
		savedSub.Username = "updated_user"

		err := suite.repo.UpdateSubscription(savedSub)
		suite.NoError(err, "Expected no error when updating subscription")

		updatedSub, err := suite.repo.FindByEmail(savedSub.Email)
		suite.NoError(err)
		suite.True(updatedSub.Verified, "Expected 'Verified' field to be updated")
		suite.Equal("updated_user", updatedSub.Username, "Expected 'Username' field to be updated")
	})
}

func (suite *SubscriptionRepoTestSuite) TestSubscriptionRepository_FindByVerificationCode() {
	suite.Run("Find by verification code", func() {
		sub := subscriptionFixture()
		verificationCode := "testVerificationCode"
		sub.VerificationCode = verificationCode
		suite.repo.SaveSubscription(sub)

		findOne, err := suite.repo.FindByVerificationCode(verificationCode)
		suite.NoError(err, "Expected no error when find subscription by verification code")

		suite.NotNil(findOne)
	})
}
