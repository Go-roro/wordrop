package main

import (
	"log"
	"net/http"
	"wordrop/internal/domain/infra/db"
	"wordrop/internal/domain/word"
)

func main() {
	database, err := db.NewMongoDatabase("mongodb://localhost:27017", "wordrop")
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	repo := db.NewWordMongoRepo(database)
	service := word.NewWordService(repo)
	r := SetupRouter(service)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
