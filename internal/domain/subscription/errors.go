package subscription

import "errors"

var (
	ErrAlreadyVerified    = errors.New("email is already verified")
	ErrRequestTooSoon     = errors.New("verification request sent too recently")
	ErrVerificationBanned = errors.New("account is banned from verification attempts")
)
