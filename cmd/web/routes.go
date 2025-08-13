package web

import (
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web/handlers"
	"github.com/Go-roro/wordrop/internal/subscription"
	"github.com/Go-roro/wordrop/internal/word"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(wordService *word.Service, subscriptionService *subscription.Service) http.Handler {
	r := chi.NewRouter()
	wordHandler := &handlers.WordHandler{WordService: wordService}
	subscriptionHandler := &handlers.SubscriptionHandler{SubscriptionService: subscriptionService}

	r.Route("/words", func(r chi.Router) {
		r.Post("/", wordHandler.SaveWordHandler)
		r.Put("/", wordHandler.UpdateWordHandler)
		r.Get("/", wordHandler.GetWordsHandler)
	})

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", subscriptionHandler.SaveNewSubscription)
	})
	return r
}
