package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/turgut-nergin/tesodev_work1/docs"

	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/broker"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/mongo"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"github.com/turgut-nergin/tesodev_work1/internal/routes"
)

// @title USER SERVICE
// @version 1.0
// @description User Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @schemes http

// @BasePath /

func main() {
	mongoConfig := config.MongoEnvConfig[os.Getenv("CURRENT_STATE")]

	config := config.RabbitMQEnvConfig["local"]
	fmt.Println(config)
	conn := broker.CreateConnection(config)
	defer conn.Close()

	ch := broker.CreateChannel(conn)
	defer ch.Close()

	q := broker.CreateQueue(ch)

	rabbitMQModel := models.RabbitMQ{
		Channel: ch,
		Queue:   q,
	}

	url := fmt.Sprintf("mongodb://%s:%s", mongoConfig.Host, mongoConfig.Port)
	client := mongo.MongoClient(url)
	collection := client.Database(mongoConfig.DBName).Collection(mongoConfig.CollectionName)
	repo := repository.New(collection)
	handler := handler.New(repo, &mongoConfig, rabbitMQModel)
	echo := echo.New()

	echo.Use(middleware.CORS())

	routes.GetRouter(echo, handler)

	log.Fatal(echo.Start(":8080"))

}
