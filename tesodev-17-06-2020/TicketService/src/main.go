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

func InitRepository(config config.Config) repository.Repositories {
	url := fmt.Sprintf("mongodb+srv://%s:%s@tesodev.4plwq.mongodb.net/?retryWrites=true&w=majority", config.UserName, config.Password)
	client := mongo.MongoClient(url)
	db := client.Database(config.DBName)
	var modelRepo repository.Repositories
	modelRepo.TicketRepository = repository.NewTicket(db.Collection(config.TicketCollection))
	modelRepo.AnswerRepository = repository.NewAnswer(db.Collection(config.AnswerCollection))
	return modelRepo
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
	routes.GetRouter(echo, handler)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error env file")
	}
	log.Fatal(echo.Start(":8080"))

}
