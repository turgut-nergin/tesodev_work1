package lib

import (
	"time"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
)

func RequestAssign(userId string, categoryId string, ticketRequest *models.TicketRequest) *models.Ticket {
	return &models.Ticket{
		Id:             uuid.New().String(),
		Subject:        ticketRequest.Subject,
		Status:         ticketRequest.Status,
		Body:           ticketRequest.Body,
		CategoryId:     categoryId,
		Attachments:    ticketRequest.Attachments,
		CreatedAt:      TimeStampNow(),
		LastAnsweredAt: TimeStampNow(),
		UpdatedAt:      TimeStampNow(),
		CreatedBy:      userId,
	}

}

func ResponseAssign(answers []models.Answer, tickets *models.Ticket) *models.TicketResponse {
	return &models.TicketResponse{
		Answers:        answers,
		Body:           tickets.Body,
		CreatedAt:      time.Unix(tickets.CreatedAt, 0),
		UpdatedAt:      time.Unix(tickets.UpdatedAt, 0),
		LastAnsweredAt: time.Unix(tickets.LastAnsweredAt, 0),
		Status:         tickets.Status,
		Attachments:    tickets.Attachments,
		Subject:        tickets.Subject,
		Id:             tickets.Id,
		CreatedBy:      tickets.CreatedBy,
	}

}
