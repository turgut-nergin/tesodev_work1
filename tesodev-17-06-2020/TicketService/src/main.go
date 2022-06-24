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
	"github.com/turgut-nergin/tesodev_work1/internal/repository/models"
	"github.com/turgut-nergin/tesodev_work1/internal/routes"
)

func InitRepository(url string, config config.Config) models.Repositorys {
	client := mongo.MongoClient(url)
	db := client.Database(config.DBName)
	var modelRepo models.Repositorys
	modelRepo.TicketR = repository.NewTicket(db)
	modelRepo.AttachmentR = repository.NewAttachment(db)
	modelRepo.AnswerR = repository.NewAnswer(db)
	return modelRepo
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("env load error")
	}
	appEnv := os.Getenv("CURRENT_STATE")
	config := config.EnvConfig[appEnv]
	url := fmt.Sprintf("mongodb+srv://%s:%s@cluster2.tw0oy.mongodb.net/?retryWrites=true&w=majority", config.UserName, config.Password)
	repositorys := InitRepository(url, config)
	handler := handler.New(repositorys)
	echo := echo.New()
	routes.GetRouter(echo, handler)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error env file")
	}
	log.Fatal(echo.Start(":8080"))

}
