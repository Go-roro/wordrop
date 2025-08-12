package web

import (
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web/handlers"
	"github.com/Go-roro/wordrop/internal/word"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(service *word.Service) http.Handler {
	r := chi.NewRouter()
	r.Route("/words", func(r chi.Router) {
		wordHandler := &handlers.WordHandler{WordService: service}
		r.Post("/", wordHandler.SaveWordHandler)
		r.Put("/", wordHandler.UpdateWordHandler)
		r.Get("/", wordHandler.GetWordsHandler)
	})
	return r
}
