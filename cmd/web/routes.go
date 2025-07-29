package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"wordrop/cmd/web/handlers"
	"wordrop/internal/domain/word"
)

func SetupRouter(service *word.Service) http.Handler {
	r := chi.NewRouter()
	r.Route("/words", func(r chi.Router) {
		wordHandler := &handlers.WordHandler{WordService: service}
		r.Post("/", wordHandler.PostWordHandler)
	})
	return r
}
