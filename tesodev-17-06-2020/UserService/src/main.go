package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
	"github.com/turgut-nergin/tesodev_work1/internal/mongo"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"github.com/turgut-nergin/tesodev_work1/internal/routes"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("env load error")
	}
	appEnv := os.Getenv("CURRENT_STATE")
	config := config.EnvConfig[appEnv]
	url := fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	client := mongo.MongoClient(url)
	collection := client.Database(config.DBName).Collection(config.CollectionName)
	repo := repository.New(collection)
	handler := handler.New(repo)
	echo := echo.New()
	routes.GetRouter(echo, handler)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error env file")
	}
	log.Fatal(echo.Start(":8080"))

}
