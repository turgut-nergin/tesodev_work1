package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/turgut-nergin/tesodev_work1/client"
	"github.com/turgut-nergin/tesodev_work1/config"
	_ "github.com/turgut-nergin/tesodev_work1/docs"
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

// @title Ticket Service
// @version 1.0
// @description Ticket Service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @schemes http

// @BasePath /
func main() {
	config := config.EnvConfig["local"]
	repositories := InitRepository(config)
	clients := GetClients()
	handler := handler.New(repositories, clients)
	echo := echo.New()
	routes.GetRouter(echo, handler)
	log.Fatal(echo.Start(":8081"))
}
