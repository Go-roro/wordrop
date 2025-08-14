package subscription

import (
	"errors"
	"fmt"

	"github.com/Go-roro/wordrop/internal/auth"
)

type Repository interface {
	FindByEmail(email string) (*Subscription, error)
	SaveSubscription(subscription *Subscription) (*Subscription, error)
	UpdateSubscription(subscription *Subscription) error
	FindByVerificationCode(code string) (*Subscription, error)
}

type MailSender interface {
	SendVerificationEmail(email, username, code string) error
}

type Service struct {
	repository  Repository
	mailSender  MailSender
	jwtProvider *auth.JwtProvider
}

func NewSubscriptionService(repo Repository, mailSender MailSender, provider *auth.JwtProvider) *Service {
	return &Service{
		repository:  repo,
		mailSender:  mailSender,
		jwtProvider: provider,
	}
}

func (s *Service) SaveSubscription(saveDto *SaveSubscriptionDto) error {
	subscription, err := s.repository.FindByEmail(saveDto.Email)
	if err != nil && errors.Is(err, ErrSubscriptionNotFound) {
		newSubscription := NewSubscription(saveDto.Username, saveDto.Email)
		subscription, err = s.repository.SaveSubscription(newSubscription)
		if err != nil {
			return fmt.Errorf("failed to save subscription: %w", err)
		}
		return s.sendVerificationEmail(subscription)
	}

	if err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	subscription.refreshBannedStatus()
	if err := subscription.validateVerifiable(); err != nil {
		return fmt.Errorf("failed to validate subscription: %w", err)
	}

	if subscription.isShouldBeBan() {
		subscription.ban()
		if err := s.repository.UpdateSubscription(subscription); err != nil {
			return fmt.Errorf("failed to update subscription after banning: %w", err)
		}
		return ErrVerificationBanned
	}

	return s.sendVerificationEmail(subscription)
}

func (s *Service) sendVerificationEmail(subscription *Subscription) error {
	subscription.refreshVerificationCode()
	token, err := s.jwtProvider.GenerateVerificationToken(subscription.ID.Hex(), subscription.VerificationCode)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	if err := s.mailSender.SendVerificationEmail(subscription.Email, subscription.Username, token); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	subscription.verificationMailSent()
	if err := s.repository.UpdateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to update subscription after sending email: %w", err)
	}
	return nil
}
