package main

import (
	"log"
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web"
	"github.com/Go-roro/wordrop/internal/infra/db"
	"github.com/Go-roro/wordrop/internal/infra/email"
	"github.com/Go-roro/wordrop/internal/subscription"
	"github.com/Go-roro/wordrop/internal/word"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

const localPort string = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := setupDatabase()
	wordRepo := word.NewWordRepo(database)
	wordService := word.NewWordService(wordRepo)

	subscriptionRepo := subscription.NewSubscriptionRepo(database)
	sender := setupMailSender()
	subscriptionService := subscription.NewSubscriptionService(subscriptionRepo, sender)

	r := web.SetupRouter(wordService, subscriptionService)
	log.Printf("Starting server on %s\n", localPort)

	if err := http.ListenAndServe(localPort, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupMailSender() *email.GmailSender {
	config, err := email.NewMailSenderConfig()
	if err != nil {
		log.Fatalf("Failed to create MailSenderConfig: %v", err)
	}
	sender, err := email.NewMailSender(config)
	if err != nil {
		log.Fatalf("Failed to create MailSender: %v", err)
	}
	return sender
}

func setupDatabase() *mongo.Database {
	database, err := db.NewMongoDatabase("mongodb://localhost:27017", "wordrop")
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}
	return database
}
