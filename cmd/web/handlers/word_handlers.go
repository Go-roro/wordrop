package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web/dto"
	"github.com/Go-roro/wordrop/internal/domain/word"
)

type WordHandler struct {
	WordService *word.Service
}

func (h *WordHandler) SaveWordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.SaveWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	saveDto := req.ToSaveDto()
	createdWord, err := h.WordService.SaveNewWord(saveDto)
	if err != nil {
		http.Error(w, "Failed to save word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdWord); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *WordHandler) UpdateWordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.UpdateWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.WordService.UpdateWord(req.ToUpdateDto())
	if err != nil {
		http.Error(w, "Failed to update word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
