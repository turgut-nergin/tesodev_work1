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
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("env load error")
	}
	appEnv := os.Getenv("CURRENT_STATE")
	config := config.EnvConfig[appEnv]
	url := fmt.Sprintf("mongodb+srv://%s:%s@cluster2.tw0oy.mongodb.net/?retryWrites=true&w=majority", config.UserName, config.Password)
	client := mongo.MongoClient(url)
	collection := client.Database(config.DBName).Collection(config.CollectionName)
	repo := repository.New(collection)
	handler := handler.New(repo)
	e := echo.New()
	e.GET("/user/:userId", handler.GetUser)
	e.POST("/user", handler.UpsertUser)
	e.DELETE("/user/:userId", handler.DeleteUser)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error env file")
	}
	log.Fatal(e.Start(":8080"))

}
