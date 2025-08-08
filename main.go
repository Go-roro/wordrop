package main

import (
	"github.com/Go-roro/wordrop/internal/domain/word"
	"github.com/Go-roro/wordrop/internal/infra/db"
	"log"
	"net/http"

	"github.com/Go-roro/wordrop/cmd/web"
	"go.mongodb.org/mongo-driver/mongo"
)

const localPort string = ":8080"

func main() {
	database := setupDatabase()
	repo := word.NewWordMongoRepo(database)
	service := word.NewWordService(repo)

	r := web.SetupRouter(service)
	log.Printf("Starting server on %s\n", localPort)

	if err := http.ListenAndServe(localPort, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupDatabase() *mongo.Database {
	database, err := db.NewMongoDatabase("mongodb://localhost:27017", "wordrop")
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}
	return database
}
