package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/turgut-nergin/tesodev_work1/client"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
)

type Handler struct {
	repository.Repositories
	clients map[string]client.Client
}

func New(repositories repository.Repositories, client map[string]client.Client) *Handler {
	return &Handler{
		Repositories: repositories,
		clients:      client,
	}
}

// TicketService
// @Summary  Create a Ticket by user and category Id
// @Description Create Ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param models.TicketRequest body models.TicketRequest true "For Create a Ticket"
// @Param ticketId query string false "Ticket ID"
// @Param userId query string false "User ID"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} string
// @Router /ticket [POST]
func (h *Handler) CreateTicket(c echo.Context) error {
	userId := c.QueryParam("userId")
	categoryId := c.QueryParam("categoryId")

	if userId == "" {
		return errors.ValidationError.WrapErrorCode(3000).WrapDesc("user id cannot be empty").ToResponse(c)
	}
	if _, err := uuid.Parse(userId); err != nil {
		return errors.ValidationError.WrapErrorCode(3001).WrapDesc(err.Error()).ToResponse(c)
	}

	ticketRequest := models.TicketRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&ticketRequest); err != nil {
		return errors.ValidationError.WrapErrorCode(3002).WrapDesc(err.Error()).ToResponse(c)
	}

	if err := lib.Validate(ticketRequest); err != nil {
		return errors.ValidationError.WrapErrorCode(2999).WrapDesc(err.Error()).ToResponse(c)

	}

	isExist, error := h.clients["userClient"].UserIsExist(userId)

	if error != nil {
		return error.ToResponse(c)

	}

	if *isExist == false {
		return errors.UnknownError.WrapErrorCode(3022).WrapDesc("user id not found!").ToResponse(c)

	}

	ticket := lib.RequestAssign(&ticketRequest)
	ticket.CreatedBy = userId
	ticket.CategoryId = categoryId

	ticketId, err := h.TicketRepository.InsertTicket(ticket)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(3003).WrapDesc(err.Error()).ToResponse(c)
	}

	return c.JSON(http.StatusCreated, ticketId)

}

// DeleteTicket by ID
// @Summary  Delete Ticket
// @Description Delete Ticket by ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /ticket/{ticketId} [delete]
func (h *Handler) DeleteTicket(c echo.Context) error {
	ticketId := c.Param("ticketId")
	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(3005).WrapDesc(err.Error()).ToResponse(c)
	}

	deletedCount, err := h.TicketRepository.DeleteTicket(ticketId)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(3006).WrapDesc(err.Error()).ToResponse(c)
	}

	if deletedCount != 0 {
		if _, err := h.AnswerRepository.DeleteAnswer(ticketId); err != nil {
			return errors.UnknownError.WrapErrorCode(3008).WrapDesc(err.Error()).ToResponse(c)
		}
		return c.JSON(http.StatusOK, true)

	}
	return c.JSON(http.StatusNotFound, false)

}

// GetTicket
// @Summary  Get Ticket by Id
// @Description Get Ticket by ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticketId query string false "Ticket ID"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.TicketResponse
// @Router /ticket/{ticketId} [GET]
func (h *Handler) GetTicket(c echo.Context) error {
	ticketId := c.Param("ticketId")

	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(3009).WrapDesc(err.Error()).ToResponse(c)
	}

	tickets, err := h.TicketRepository.GetTicket(ticketId)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(3010).WrapDesc(err.Error()).ToResponse(c)
	}

	if tickets == nil {
		return errors.NotFound.WrapErrorCode(11000).
			WrapDesc(fmt.Sprintf("Ticket id: %v not found", ticketId)).ToResponse(c)
	}

	category, error := h.clients["categoryClient"].GetCategory(tickets.CategoryId)

	if error != nil {
		return error.ToResponse(c)
	}

	answers, err := h.AnswerRepository.GetAnswers(ticketId)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(3011).WrapDesc(err.Error()).ToResponse(c)
	}

	ticketsResponse := lib.TicketResponseAssign(tickets)

	ticketsResponse.Answers = lib.AnswersAssign(answers)

	ticketsResponse.Category = category

	return c.JSON(http.StatusNotFound, ticketsResponse)

}

// UpdateAnswer by ID
// @Summary  Update Answer
// @Description Update Answer by ID
// @Tags Answers
// @Accept json
// @Produce json
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /ticket/{ticketId} [PUT]
func (h *Handler) UpdateAnswer(c echo.Context) error {
	answerId := c.Param("answerId")
	if _, err := uuid.Parse(answerId); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	answer := models.Answer{}

	if err := json.NewDecoder(c.Request().Body).Decode(&answer); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	answer.UpdatedAt = lib.TimeStampNow()
	answer.Id = answerId
	modifiedCount, err := h.AnswerRepository.UpdateAnswer(&answer)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(2027).WrapDesc(err.Error()).ToResponse(c)

	}

	if *modifiedCount == 0 {
		return c.JSON(http.StatusNotFound, false)

	}

	return c.JSON(http.StatusOK, true)

}

// CreateAnswer
// @Summary  Create Answer
// @Description Create Answer by user and ticket id
// @Tags Answers
// @Accept json
// @Produce json
// @Param userId query string false "User ID"
// @Param models.Answer body models.Answer true "For Create an Answer"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} string
// @Router /ticket/{ticketId}/answer [POST]
func (h *Handler) CreateAnswer(c echo.Context) error {
	ticketId := c.Param("ticketId")
	userId := c.QueryParam("userId")

	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(2029).WrapDesc(err.Error()).ToResponse(c)
	}

	if _, err := uuid.Parse(userId); err != nil {
		return errors.ValidationError.WrapErrorCode(2030).WrapDesc(err.Error()).ToResponse(c)
	}

	answer := models.Answer{}

	if err := json.NewDecoder(c.Request().Body).Decode(&answer); err != nil {
		return errors.ValidationError.WrapErrorCode(2031).WrapDesc(err.Error()).ToResponse(c)
	}

	userIsExist, error := h.clients["userClient"].UserIsExist(userId)

	if error != nil {
		return error.ToResponse(c)

	}

	if *userIsExist == false {
		return errors.UnknownError.WrapErrorCode(2032).WrapDesc("user id not found!").ToResponse(c)

	}
	answer.CreatedAt = lib.TimeStampNow()
	answer.TicketId = ticketId
	answer.Id = uuid.New().String()
	answer.UserId = userId

	answerId, err := h.AnswerRepository.CreateAnswer(&answer)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(2033).WrapDesc(err.Error()).ToResponse(c)
	}

	return c.JSON(http.StatusCreated, answerId)

}
