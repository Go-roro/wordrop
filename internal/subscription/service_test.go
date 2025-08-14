package subscription

import (
	"github.com/Go-roro/wordrop/internal/auth"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SubscriptionServiceTestSuite struct {
	suite.Suite
	mockRepo       *MockRepository
	mockMailSender *MockMailSender
	service        *Service
}

func (suite *SubscriptionServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.mockMailSender = new(MockMailSender)
	provider, _ := auth.NewJwtProvider("a-string-secret-at-least-256-bits-long")
	suite.service = NewSubscriptionService(suite.mockRepo, suite.mockMailSender, provider)
}

func TestSubscriptionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionServiceTestSuite))
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_NewUser_Success() {
	// Given
	dto := &SaveSubscriptionDto{Email: "new@example.com", Username: "NewUser"}
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(nil, ErrSubscriptionNotFound)
	suite.mockRepo.EXPECT().SaveSubscription(mock.AnythingOfType("*subscription.Subscription")).Return(
		&Subscription{
			Email:    dto.Email,
			Username: dto.Username,
		}, nil)
	suite.mockMailSender.EXPECT().SendVerificationEmail(dto.Email, dto.Username, mock.AnythingOfType("string")).Return(nil)
	suite.mockRepo.EXPECT().UpdateSubscription(mock.AnythingOfType("*subscription.Subscription")).Return(nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_ExistingUser_Success() {
	// Given
	dto := &SaveSubscriptionDto{Email: "exist@example.com", Username: "ExistUser"}
	existingSub := NewSubscription(dto.Username, dto.Email)
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(existingSub, nil)
	suite.mockMailSender.EXPECT().SendVerificationEmail(
		existingSub.Email,
		existingSub.Username,
		mock.AnythingOfType("string"),
	).Return(nil)
	suite.mockRepo.EXPECT().UpdateSubscription(existingSub).Return(nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertExpectations(suite.T())
	suite.Equal(existingSub.VerificationAttempts, 1)
	suite.NotNil(existingSub.VerificationCode)
	suite.NotNil(existingSub.LastVerifiedAt)
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_BannedExpiredUser_Success() {
	// Given
	dto := &SaveSubscriptionDto{Email: "user@example.com", Username: "user"}
	sub := NewSubscription(dto.Username, dto.Email)
	sub.Banned = true
	sub.BannedUntil = time.Now() // Expired ban
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(sub, nil)
	suite.mockMailSender.EXPECT().SendVerificationEmail(
		sub.Email,
		sub.Username,
		mock.AnythingOfType("string"),
	).Return(nil)
	suite.mockRepo.EXPECT().UpdateSubscription(sub).Return(nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertExpectations(suite.T())
	suite.False(sub.Banned)
	//suite.Nil(sub.BannedUntil)

	suite.Equal(sub.VerificationAttempts, 1)
	suite.NotNil(sub.VerificationCode)
	suite.NotNil(sub.LastVerifiedAt)
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_BannedUser() {
	// Given
	dto := &SaveSubscriptionDto{Email: "banned@example.com", Username: "BannedUser"}
	bannedSub := NewSubscription(dto.Username, dto.Email)
	bannedSub.Banned = true
	bannedSub.BannedUntil = time.Now().Add(24 * time.Hour)
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(bannedSub, nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.ErrorIs(err, ErrVerificationBanned)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertNotCalled(suite.T(), "SendVerificationEmail")
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_TooRapidRequestedUser() {
	// Given
	dto := &SaveSubscriptionDto{Email: "user@example.com", Username: "user"}
	user := NewSubscription(dto.Username, dto.Email)
	user.LastVerifiedAt = time.Now() // Assuming user has just verified
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(user, nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.ErrorIs(err, ErrRequestTooSoon)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertNotCalled(suite.T(), "SendVerificationEmail")
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_TooManyRequestedUser() {
	// Given
	dto := &SaveSubscriptionDto{Email: "user@example.com", Username: "user"}
	user := NewSubscription(dto.Username, dto.Email)
	user.VerificationAttempts = maxVerificationAttempts
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(user, nil)
	suite.mockRepo.EXPECT().UpdateSubscription(user).Return(nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.ErrorIs(err, ErrVerificationBanned)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertNotCalled(suite.T(), "SendVerificationEmail")

	suite.True(user.Banned)
	suite.NotNil(user.BannedUntil)
}

func (suite *SubscriptionServiceTestSuite) TestSaveSubscription_AlreadyVerifiedUser() {
	// Given
	dto := &SaveSubscriptionDto{Email: "user@example.com", Username: "user"}
	user := NewSubscription(dto.Username, dto.Email)
	user.Verified = true
	suite.mockRepo.EXPECT().FindByEmail(dto.Email).Return(user, nil)

	// When
	err := suite.service.SaveSubscription(dto)

	// Then
	suite.ErrorIs(err, ErrAlreadyVerified)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockMailSender.AssertNotCalled(suite.T(), "SendVerificationEmail")
}
