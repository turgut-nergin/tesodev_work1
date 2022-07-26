package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func SetGeneric(echo *echo.Echo, handler *handler.Handler) {
	echo.GET("/swagger/*", echoSwagger.WrapHandler)
	echo.POST("/csv", handler.ReadCsv)

}

func SetAdminRouter(echo *echo.Group, handler *handler.Handler) {
	echo.DELETE("/category/:categoryId", handler.DeleteCategory)
	echo.POST("/category", handler.CreateCategory)
}

func SetUserRouter(echo *echo.Group, handler *handler.Handler) {
	echo.GET("/category", handler.GetCategory)
	echo.PUT("/category/:categoryId", handler.UpdateCategory)
	echo.GET("/categories", handler.GetCategories)
}
