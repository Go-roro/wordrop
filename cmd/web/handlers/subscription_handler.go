package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web/dto"
	"github.com/Go-roro/wordrop/internal/subscription"
)

type SubscriptionHandler struct {
	SubscriptionService *subscription.Service
}

func (h *SubscriptionHandler) SaveNewSubscription(w http.ResponseWriter, r *http.Request) {
	var req *dto.SaveSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		NewHTTPError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	saveDto := req.ToSaveDto()
	err := h.SubscriptionService.SaveSubscription(saveDto)
	if err != nil {
		errMessage := fmt.Errorf("failed to save subscription: %v", err)
		NewHTTPError(w, errMessage.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
