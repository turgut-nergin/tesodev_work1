package routes

import (
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/client"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler, clients map[string]client.Client) {

	echo.POST("/ticket", handler.CreateTicket(clients["userClient"]))
	echo.POST("/ticket/:ticketId/answer", handler.CreateAnswer(clients["userClient"]))
	echo.PUT("/ticket/:answerId/answer", handler.UpdateAnswer)
	echo.GET("/ticket/:ticketId", handler.GetTicket)
	echo.DELETE("/ticket/:ticketId", handler.DeleteTicket)

}
