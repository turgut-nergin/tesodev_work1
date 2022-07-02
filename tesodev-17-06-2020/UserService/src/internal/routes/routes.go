package routes

import (
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler) {
	echo.GET("/user/:userId", handler.GetUser)
	echo.POST("/user", handler.UpsertUser)
	echo.DELETE("/user/:userId", handler.DeleteUser)
}
