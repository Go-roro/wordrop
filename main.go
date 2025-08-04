package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"wordrop/cmd/web"
	"wordrop/internal/domain/infra/db"
	"wordrop/internal/domain/word"
)

const localPort string = ":8080"

func main() {
	database := setupDatabase()
	repo := db.NewWordMongoRepo(database)
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
