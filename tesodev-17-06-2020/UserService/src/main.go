package main

import (
	"log"

	"github.com/turgut-nergin/tesodev_work1/internal/handler"
	"github.com/turgut-nergin/tesodev_work1/internal/mongo"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"github.com/turgut-nergin/tesodev_work1/internal/router"
)

func main() {
	url := "mongodb://mongo-db:27017"
	client := mongo.MongoClient(url)
	collection := client.Database("Tickets").Collection("users")
	repo := repository.New(collection)
	handler := handler.New(repo)
	e := router.New(handler)
	log.Fatal(e.Start(":8080"))

}
