package handlers

import (
	"encoding/json"
	"net/http"
	"wordrop/cmd/web/dto"
	"wordrop/internal/domain/word"
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
