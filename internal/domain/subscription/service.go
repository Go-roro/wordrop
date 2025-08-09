package subscription

import (
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
	if s.isNewEmail(saveDto.Email) {
		subscription := NewSubscription(saveDto.Email, saveDto.Username)
		savedSubscription, err := s.repository.SaveSubscription(subscription)
		if err != nil {
			return fmt.Errorf("failed to save subscription: %w", err)
		}
		return s.sendVerificationEmail(savedSubscription)
	}

	existSubscription, err := s.repository.FindByEmail(saveDto.Email)
	if err != nil {
		return fmt.Errorf("failed to find existing subscription: %w", err)
	}

	existSubscription.refreshBannedStatus()
	if err := existSubscription.ValidateVerifiable(); err != nil {
		return fmt.Errorf("failed to validate subscription: %w", err)
	}

	if existSubscription.isShouldBeBan() {
		existSubscription.ban()
		if err := s.repository.UpdateSubscription(existSubscription); err != nil {
			return fmt.Errorf("failed to update subscription after banning: %w", err)
		}
		return ErrVerificationBanned
	}

	return s.sendVerificationEmail(existSubscription)
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

func (s *Service) isNewEmail(email string) bool {
	return !s.repository.ExistsByEmail(email)
}
