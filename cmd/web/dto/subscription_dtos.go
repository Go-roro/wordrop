package dto

import "github.com/Go-roro/wordrop/internal/subscription"

type SaveSubscriptionRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

func (r *SaveSubscriptionRequest) ToSaveDto() *subscription.SaveSubscriptionDto {
	return &subscription.SaveSubscriptionDto{
		Email:    r.Email,
		Username: r.Username,
	}
}
