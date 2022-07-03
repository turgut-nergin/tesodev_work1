package lib

import (
	"time"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
)

func RequestAssign(ticketRequest *models.TicketRequest) *models.Ticket {
	return &models.Ticket{
		Id:             uuid.New().String(),
		Subject:        ticketRequest.Subject,
		Status:         ticketRequest.Status,
		Body:           ticketRequest.Body,
		Attachments:    ticketRequest.Attachments,
		CreatedAt:      TimeStampNow(),
		LastAnsweredAt: TimeStampNow(),
		UpdatedAt:      TimeStampNow(),
	}

}

func TicketResponseAssign(tickets *models.Ticket) *models.TicketResponse {
	return &models.TicketResponse{
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

func AnswersAssign(answers []models.Answer) []models.AnswerResponse {
	answersResponse := []models.AnswerResponse{}

	for _, answer := range answers {
		answerResponse := models.AnswerResponse{
			CreatedAt: time.Unix(answer.CreatedAt, 0),
			UpdatedAt: time.Unix(answer.UpdatedAt, 0),
			Body:      answer.Body,
			TicketId:  answer.TicketId,
			UserId:    answer.UserId,
			Id:        answer.Id,
		}
		answersResponse = append(answersResponse, answerResponse)
	}
	return answersResponse
}
