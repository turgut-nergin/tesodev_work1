package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler) {

	echo.POST("/ticket", handler.CreateTicket)
	echo.POST("/ticket/:ticketId/answer", handler.CreateAnswer)
	echo.PUT("/ticket/answer/:answerId", handler.UpdateAnswer)
	echo.GET("/ticket/:ticketId", handler.GetTicket)
	echo.DELETE("/ticket/:ticketId", handler.DeleteTicket)
	echo.GET("/ticket", handler.GetTickets)
	echo.GET("/swagger/*", echoSwagger.WrapHandler)

}
