package routes

import (
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/internal/handler"
	"github.com/turgut-nergin/tesodev_work1/pkg/user"
)

func GetRouter(echo *echo.Echo, handler *handler.Handler, clients map[string]user.Client) {

	echo.POST("/ticket", handler.CreateTicket(clients["userClient"]))
	echo.POST("/ticket/:ticketId/answer", handler.UpsertAnswer)
	echo.GET("/ticket/:ticketId", handler.GetTicket)
	echo.DELETE("/ticket/:ticketId", handler.DeleteTicket)

}
