package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"wordrop/cmd/web/handlers"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/words", func(r chi.Router) {
		r.Post("/", handlers.WordHandler.PostWordHandler)
	})
	return r
}
