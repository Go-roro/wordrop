package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Go-roro/wordrop/cmd/web/dto"
	"github.com/Go-roro/wordrop/internal/domain/word"
)

type WordHandler struct {
	WordService *word.Service
}

func (h *WordHandler) SaveWordHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.SaveWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		NewHTTPError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	saveDto := req.ToSaveDto()
	createdWord, err := h.WordService.SaveNewWord(saveDto)
	if err != nil {
		NewHTTPError(w, "Failed to save word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdWord); err != nil {
		NewHTTPError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *WordHandler) UpdateWordHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		NewHTTPError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.WordService.UpdateWord(req.ToUpdateDto())
	if err != nil {
		NewHTTPError(w, "Failed to update word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WordHandler) GetWordsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	pageSize, err := strconv.Atoi(q.Get("page_size"))
	sortBy := q.Get("sort_by")
	sortOrder := q.Get("sort_order")
	isDelivered, err := strconv.ParseBool(q.Get("is_delivered"))
	if err != nil {
		NewHTTPError(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	params := &word.SearchParams{
		Page:        page,
		PageSize:    pageSize,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
		IsDelivered: &isDelivered,
	}

	words, err := h.WordService.FindWords(params)
	if err != nil {
		NewHTTPError(w, "Failed to retrieve words", http.StatusInternalServerError)
		return
	}

	if page > words.LastPage {
		NewHTTPError(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(words); err != nil {
		NewHTTPError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
