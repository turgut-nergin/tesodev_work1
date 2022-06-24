package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	repositoryModel "github.com/turgut-nergin/tesodev_work1/internal/repository/models"
)

type Handler struct {
	repositoryModel.Repositorys
}

func New(repositorys repositoryModel.Repositorys) *Handler {
	return &Handler{
		Repositorys: repositorys,
	}
}
func (h *Handler) CreateTicket(c echo.Context) error {
	userid := c.QueryParam("userId")
	categoryId := c.QueryParam("categoryId")

	if userid == "" {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc("user id cannot be empty").ToResponse(c)
	}
	if _, err := uuid.Parse(userid); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	ticketRequest := models.TicketRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&ticketRequest); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	ticket := models.Ticket{}
	ticket.Id = uuid.New().String()
	ticket.Subject = ticketRequest.Subject
	ticket.Status = ticketRequest.Status
	ticket.Body = ticketRequest.Body
	ticket.CategoryId = categoryId
	ticket.CreatedAt = lib.TimeStampNow()
	ticket.LastAnsweredAt = lib.TimeStampNow()
	ticket.UpdatedAt = lib.TimeStampNow()
	ticket.CreatedBy = userid

	ticketId, err := h.TicketR.InsertTicket(&ticket)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(2001).WrapDesc(err.Error()).ToResponse(c)
	}

	err = h.AttachmentR.InsertAttachment(*ticketId, ticketRequest.Attachments)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(2002).WrapDesc(err.Error()).ToResponse(c)
	}

	return c.JSON(http.StatusCreated, ticketId)

}

func (h *Handler) DeleteTicket(c echo.Context) error {
	ticketId := c.Param("ticketId")
	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	deletedCount, err := h.TicketR.DeleteTicket(ticketId)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(2002).WrapDesc(err.Error()).ToResponse(c)
	}

	if deletedCount != 0 {
		if _, err := h.AttachmentR.DeleteAttachment(ticketId); err != nil {
			return errors.UnknownError.WrapErrorCode(2002).WrapDesc(err.Error()).ToResponse(c)
		}

		if _, err := h.AnswerR.DeleteAnswer(ticketId); err != nil {
			return errors.UnknownError.WrapErrorCode(2002).WrapDesc(err.Error()).ToResponse(c)
		}
		return c.JSON(http.StatusOK, true)

	}
	return c.JSON(http.StatusNotFound, false)

}

func (h *Handler) GetTicket(c echo.Context) error {
	ticketId := c.Param("ticketId")

	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	tickets, err := h.TicketR.GetTicket(ticketId)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(2005).WrapDesc(err.Error()).ToResponse(c)
	}

	if tickets == nil {
		return errors.NotFound.WrapErrorCode(11000).
			WrapDesc(fmt.Sprintf("Ticket id: %v not found", ticketId)).ToResponse(c)
	}

	answers, err := h.AnswerR.GetAnswers(ticketId)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(2004).WrapDesc(err.Error()).ToResponse(c)
	}

	attachments, err := h.AttachmentR.GetAttachment(ticketId)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(2002).WrapDesc(err.Error()).ToResponse(c)
	}

	ticketsResponse := models.TicketResponse{}
	ticketsResponse.Answers = answers
	ticketsResponse.Attachments = attachments
	ticketsResponse.Body = tickets.Body
	ticketsResponse.CreatedAt = time.Unix(tickets.CreatedAt, 0)
	ticketsResponse.UpdatedAt = time.Unix(tickets.UpdatedAt, 0)
	ticketsResponse.LastAnsweredAt = time.Unix(tickets.LastAnsweredAt, 0)
	ticketsResponse.Status = tickets.Status
	ticketsResponse.Subject = tickets.Subject
	ticketsResponse.Id = tickets.Id
	ticketsResponse.CreatedBy = tickets.CreatedBy

	return c.JSON(http.StatusNotFound, ticketsResponse)

}

func (h *Handler) UpsertAnswer(c echo.Context) error {
	ticketId := c.Param("ticketId")
	answerId := c.QueryParam("answerId")

	if answerId != "" {
		if _, err := uuid.Parse(answerId); err != nil {
			return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
		}
	} else {
		answerId = uuid.New().String()
	}

	if _, err := uuid.Parse(ticketId); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	answer := models.Answer{}

	if err := json.NewDecoder(c.Request().Body).Decode(&answer); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	answer.UpdatedAt = lib.TimeStampNow()
	answer.CreatedBy = ticketId

	modifiedCount, err := h.AnswerR.UpsertAnswer(answerId, &answer)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(2001).WrapDesc(err.Error()).ToResponse(c)
	}
	if *modifiedCount == 0 {
		return c.JSON(http.StatusCreated, true)

	}
	return c.JSON(http.StatusOK, true)

}
