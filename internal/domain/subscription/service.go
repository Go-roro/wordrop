package subscription

import (
	"errors"
	"fmt"

	"github.com/Go-roro/wordrop/internal/infra/email"
)

type Service struct {
	repository *Repository
	mailSender email.MailSender
}

func NewSubscriptionService(repo *Repository, mailSender email.MailSender) *Service {
	return &Service{
		repository: repo,
		mailSender: mailSender,
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
	if err := subscription.ValidateVerifiable(); err != nil {
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
	err := s.mailSender.SendVerificationEmail(subscription.Email, subscription.Username, subscription.VerificationCode)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	subscription.verificationMailSent()
	if err := s.repository.UpdateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to update subscription after sending email: %w", err)
	}
	return nil
}
