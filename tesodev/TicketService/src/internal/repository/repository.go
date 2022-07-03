package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	collection *mongo.Collection
}

type Repositories struct {
	TicketRepository TicketRepository
	AnswerRepository AnswerRepository
}
