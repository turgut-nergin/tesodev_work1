package router

import (
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

type Router struct {
	handler *handler.Handler
}

func New(handler *handler.Handler) *echo.Echo {
	e := echo.New()
	e.GET("/user/:userId", handler.GetUser)
	e.POST("/user", handler.UpsertUser)
	e.DELETE("/user/:userId", handler.DeleteUser)
	return e
}
