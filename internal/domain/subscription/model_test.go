package subscription

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscription_ValidateVerifiable(t *testing.T) {
	tests := []struct {
		name    string
		sub     *Subscription
		wantErr error
	}{
		{
			name: "Already Verified",
			sub: &Subscription{
				Verified: true,
			},
			wantErr: ErrAlreadyVerified,
		},
		{
			name: "Request Too Soon",
			sub: &Subscription{
				Verified:       false,
				LastVerifiedAt: time.Now().Add(-verificationCooldown / 2), // less than cooldown period
			},
			wantErr: ErrRequestTooSoon,
		},
		{
			name: "Banned User",
			sub: &Subscription{
				Verified:    false,
				Banned:      true,
				BannedUntil: time.Now().Add(24 * time.Hour),
			},
			wantErr: ErrVerificationBanned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sub.ValidateVerifiable()
			assert.Error(t, err)
		})
	}
}

func TestSubscription_RefreshBannedStatus(t *testing.T) {
	tests := []struct {
		name         string
		subscription *Subscription
		banned       bool
	}{
		{
			name: "Banned User Still Banned",
			subscription: &Subscription{
				Banned:      true,
				BannedUntil: time.Now().Add(24 * time.Hour), // still banned
			},
			banned: true,
		},
		{
			name: "Renew Banned User After ban Period",
			subscription: &Subscription{
				Banned:      true,
				BannedUntil: time.Now().Add(-24 * time.Hour), // ban period expired
			},
			banned: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.subscription.refreshBannedStatus()
			assert.Equal(t, tt.banned, tt.subscription.Banned)
		})
	}
}

func TestSubscription_isShouldBeBan(t *testing.T) {
	tests := []struct {
		name         string
		subscription *Subscription
		shouldBeBan  bool
	}{
		{
			name:         "Should Be Banned",
			subscription: &Subscription{VerificationAttempts: maxVerificationAttempts},
			shouldBeBan:  true,
		},
		{
			name:         "Should Not Be Banned",
			subscription: &Subscription{VerificationAttempts: maxVerificationAttempts - 1},
			shouldBeBan:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.subscription.isShouldBeBan()
			assert.Equal(t, tt.shouldBeBan, got)
		})
	}
}

func TestSubscription_ban(t *testing.T) {
	tests := []struct {
		name         string
		subscription *Subscription
	}{
		{
			name: "Ban User",
			subscription: &Subscription{
				Banned:               false,
				VerificationAttempts: maxVerificationAttempts,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.subscription.ban()
			assert.True(t, tt.subscription.Banned)
			assert.NotEqual(t, time.Time{}, tt.subscription.BannedUntil)
			assert.True(t, tt.subscription.BannedUntil.After(time.Now()))
		})
	}
}

func TestSubscription_RefreshVerificationCode(t *testing.T) {
	tests := []struct {
		name         string
		subscription *Subscription
	}{
		{
			name: "Refresh Verification Code",
			subscription: &Subscription{
				VerificationCode: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.subscription.refreshVerificationCode()
			assert.NotEmpty(t, tt.subscription.VerificationCode)
			assert.Equal(t, verificationCodeLength, len(tt.subscription.VerificationCode))
		})
	}
}

func TestSubscription_VerificationMailSent(t *testing.T) {
	tests := []struct {
		name         string
		subscription *Subscription
	}{
		{
			name: "Mark Verification Mail as Sent",
			subscription: &Subscription{
				VerificationCode:     "",
				VerificationAttempts: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.subscription.verificationMailSent()
			assert.Equal(t, 1, tt.subscription.VerificationAttempts)
			assert.NotEqual(t, time.Time{}, tt.subscription.LastVerifiedAt)
		})
	}
}
