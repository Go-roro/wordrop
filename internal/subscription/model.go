package subscription

import (
	"time"

	"github.com/Go-roro/wordrop/internal/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	maxVerificationAttempts = 3
	verificationCodeLength  = 24
	verificationCooldown    = 1 * time.Minute
	banDuration             = 24 * time.Hour
)

type Subscription struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Username             string             `bson:"username" validate:"required"`
	Email                string             `bson:"email" validate:"required,email"`
	Verified             bool               `bson:"verified"`
	VerificationAttempts int                `bson:"verification_attempts"`
	LastVerifiedAt       time.Time          `bson:"last_verified_at"`
	Banned               bool               `bson:"banned"`
	BannedUntil          time.Time          `bson:"banned_until"`
	VerificationCode     string             `bson:"verification_code"`
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
}

func NewSubscription(username, email string) *Subscription {
	return &Subscription{
		Username:             username,
		Email:                email,
		Verified:             false,
		VerificationAttempts: 0,
		LastVerifiedAt:       time.Time{},
		Banned:               false,
		BannedUntil:          time.Time{},
		VerificationCode:     "",
	}
}

func (s *Subscription) validateVerifiable() error {
	if s.Verified && !s.Banned {
		return ErrAlreadyVerified
	}

	if time.Since(s.LastVerifiedAt) < verificationCooldown {
		return ErrRequestTooSoon
	}

	if s.Banned && time.Now().Before(s.BannedUntil) {
		return ErrVerificationBanned
	}
	return nil
}

func (s *Subscription) refreshBannedStatus() {
	if s.Banned && s.BannedUntil.Before(time.Now()) {
		s.Banned = false
		s.VerificationAttempts = 0
		s.BannedUntil = time.Time{}
	}
}

func (s *Subscription) isShouldBeBan() bool {
	return s.VerificationAttempts >= maxVerificationAttempts
}

func (s *Subscription) ban() {
	s.Banned = true
	s.BannedUntil = time.Now().Add(banDuration)
	s.VerificationAttempts = 0
	s.Verified = false
}

func (s *Subscription) refreshVerificationCode() {
	token, _ := common.GenerateRandomToken(verificationCodeLength)
	s.VerificationCode = token
}

func (s *Subscription) verificationMailSent() {
	s.VerificationAttempts++
	s.LastVerifiedAt = time.Now()
}
