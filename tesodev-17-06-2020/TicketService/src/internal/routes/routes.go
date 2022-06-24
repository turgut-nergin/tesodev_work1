package routes

import (
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler) {
	echo.POST("/ticket", handler.CreateTicket)
	echo.POST("/ticket/:ticketId/answer", handler.UpsertAnswer)
	echo.GET("/ticket/:ticketId", handler.GetTicket)
	echo.DELETE("/ticket/:ticketId", handler.DeleteTicket)

}
