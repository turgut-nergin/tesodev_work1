package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/turgut-nergin/tesodev_work1/client"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Handler struct {
	repository.Repositories
	clients map[string]client.Client
	cfg     *config.Config
}

func New(repositories repository.Repositories, client map[string]client.Client, config *config.Config) *Handler {
	return &Handler{
		Repositories: repositories,
		clients:      client,
		cfg:          config,
	}
}

// TicketService
// @Summary  Create a Ticket by user and category Id
// @Description Create Ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param models.TicketRequest body models.TicketRequest true "For Create a Ticket"
// @Param categoryId query string false "Category ID"
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

	if categoryId == "" {
		return errors.ValidationError.WrapErrorCode(3000).WrapDesc("category id cannot be empty").ToResponse(c)
	}
	if _, err := uuid.Parse(categoryId); err != nil {
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

	_, err := h.clients["categoryClient"].GetCategory(categoryId)
	fmt.Println(err)
	if err != nil {
		return err.ToResponse(c)
	}

	if *isExist == false {
		return errors.UnknownError.WrapErrorCode(3022).WrapDesc("user id not found!").ToResponse(c)
	}

	ticket := lib.RequestAssign(&ticketRequest)
	ticket.CreatedBy = userId
	ticket.CategoryId = categoryId

	ticketId, err := h.TicketRepository.InsertTicket(ticket)
	if err != nil {
		return err.ToResponse(c)
	}

	return c.JSON(http.StatusCreated, ticketId)

}

// DeleteTicket by ID
// @Summary  Delete Ticket
// @Description Delete Ticket by ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticketId path string true "Ticket Id"
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
// @Param ticketId path string true "Ticket Id"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.TicketResponse
// @Router /ticket/{ticketId} [GET]
func (h *Handler) GetTicket(c echo.Context) error {
	ticketId := c.Param("ticketId")

	fmt.Println(ticketId)
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

	answers, error := h.AnswerRepository.GetAnswers(ticketId)

	if error != nil {
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
// @Param answerId path string true "Answer Id"
// @Param models.AnswerRequest body models.AnswerRequest true "For update a answer"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /ticket/answer/{answerId} [PUT]
func (h *Handler) UpdateAnswer(c echo.Context) error {
	answerId := c.Param("answerId")
	if _, err := uuid.Parse(answerId); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	answerRequest := models.AnswerRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&answerRequest); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	answer := models.Answer{}
	answer.UpdatedAt = lib.TimeStampNow()
	answer.Id = answerId
	answer.Body = answerRequest.Body
	modifiedCount, err := h.AnswerRepository.UpdateAnswer(&answer)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(2027).WrapDesc(err.Error()).ToResponse(c)

	}

	if *modifiedCount == 0 {
		return c.JSON(http.StatusNotFound, false)

	}

	fmt.Println(answerId)

	answerR, error := h.AnswerRepository.GetAnswer(answerId)

	if error != nil {
		return error.ToResponse(c)
	}
	ticket := models.Ticket{}
	ticket.Id = answerR.TicketId
	ticket.LastAnsweredAt = lib.TimeStampNow()

	modifiedCount, error = h.TicketRepository.Update(ticket)

	if error != nil {
		return error.ToResponse(c)
	}

	if *modifiedCount == 0 {
		fmt.Println("here")
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
// @Param ticketId path string true "Ticket Id"
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

// GetTickets
// @Summary  Get Tickets by params
// @Description Get Tickets by params
// @Tags Tickets
// @Accept json
// @Produce json
// @Param subject query string false "subject"
// @Param body query string false "body"
// @Param status query string false "status"
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param sort query string false "sort"
// @Param direction query string false "direction"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.CategoryResponse
// @Router /tickets [get]
func (h *Handler) GetTickets(c echo.Context) error {

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, offset := lib.ValidatePaginator(limitStr, offsetStr, h.cfg.MaxPageLimit)

	filter := map[string]interface{}{}

	if subject := c.QueryParam("subject"); subject != "" {
		filter["subject"] = bson.M{"$regex": primitive.Regex{
			Pattern: subject,
			Options: "i",
		}}
	}

	if body := c.QueryParam("body"); body != "" {
		filter["body"] = bson.M{"$regex": primitive.Regex{
			Pattern: body,
			Options: "i",
		}}
	}

	if status := c.QueryParam("status"); status != "" {
		filter["status"] = bson.M{"$regex": primitive.Regex{
			Pattern: status,
			Options: "i",
		}}
	}

	lastAnsweredAt := c.QueryParam("lastAnsweredAt")

	if lastAnsweredAt != "" {
		date, error := time.Parse("2006-01-02", lastAnsweredAt)

		if error != nil {
			return errors.ValidationError.WrapErrorCode(4050).WrapDesc("last lastAnsweredAt parse error").ToResponse(c)
		}
		l, _ := time.LoadLocation("Europe/Istanbul")
		timeStamp := date.In(l).Unix()

		filter["lastAnsweredAt"] = bson.M{"$gte": timeStamp}

	}

	sortField := lib.GetAcceptedSortField(c.QueryParam("sort"))              //for example name
	sortDirection := lib.GetAcceptedSortDirection(c.QueryParam("direction")) //asc desc -> 0,-1
	tickets, err := h.TicketRepository.Find(limit, offset, filter, sortField, sortDirection)
	if err != nil {
		return err.ToResponse(c)
	}

	return c.JSON(http.StatusOK, tickets)
}
