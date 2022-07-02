package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/client"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
	"github.com/turgut-nergin/tesodev_work1/internal/mongo"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"github.com/turgut-nergin/tesodev_work1/internal/routes"
)

func InitRepository(config config.Config) repository.Repositories {
	url := fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	client := mongo.MongoClient(url)
	db := client.Database(config.DBName)
	var modelRepo repository.Repositories
	modelRepo.TicketRepository = repository.NewTicket(db.Collection(config.TicketCollection))
	modelRepo.AnswerRepository = repository.NewAnswer(db.Collection(config.AnswerCollection))
	return modelRepo
}

func GetClients() map[string]client.Client {
	return map[string]client.Client{
		"userClient":     *client.NewClient("http://localhost:8080"),
		"categoryClient": *client.NewClient("http://localhost:8082"),
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("env load error")
	}
	appEnv := os.Getenv("CURRENT_STATE")
	config := config.EnvConfig[appEnv]
	repositories := InitRepository(config)
	handler := handler.New(repositories)
	echo := echo.New()
	clients := GetClients()
	routes.GetRouter(echo, handler, clients)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error env file")
	}
	log.Fatal(echo.Start(":8081"))
}
