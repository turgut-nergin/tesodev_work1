package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler) {
	echo.GET("/user/:userId", handler.GetUser)
	echo.POST("/user", handler.UpsertUser)
	echo.DELETE("/user/:userId", handler.DeleteUser)
	echo.GET("/user", handler.Validate)
	echo.GET("/users", handler.GetUsers)
	echo.POST("/login", handler.Login)
	echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
