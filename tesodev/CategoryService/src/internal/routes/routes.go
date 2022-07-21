package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler) {
	echo.GET("/category/:categoryId", handler.GetCategory)
	echo.POST("/category", handler.CreateCategory)
	echo.PUT("/category/:categoryId", handler.UpdateCategory)
	echo.DELETE("/category/:categoryId", handler.DeleteCategory)
	echo.GET("/category", handler.GetCategories)
	echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
