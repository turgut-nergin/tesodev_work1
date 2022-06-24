package models

import "github.com/turgut-nergin/tesodev_work1/internal/repository"

type Repositorys struct {
	TicketR     repository.TicketRepository
	AttachmentR repository.AttachmentRepository
	AnswerR     repository.AnswerRepository
}
