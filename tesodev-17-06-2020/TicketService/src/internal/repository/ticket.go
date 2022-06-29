package repository

import (
	"context"
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type TicketRepository interface {
	InsertTicket(ticket *models.Ticket) (*string, error)
	DeleteTicket(id string) (int64, error)
	GetTicket(id string) (*models.Ticket, error)
}

func NewTicket(mongoClient *mongo.Collection) *Repository {
	return &Repository{mongoClient}
}

func (r *Repository) InsertTicket(ticket *models.Ticket) (*string, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if _, err := r.collection.InsertOne(context, ticket); err != nil {
		return nil, err
	}

	return &ticket.Id, nil
}

func (r *Repository) GetTicket(id string) (*models.Ticket, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ticket := models.Ticket{}
	if err := r.collection.FindOne(context, bson.M{"_id": id}).Decode(&ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) DeleteTicket(id string) (int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	result, err := r.collection.DeleteOne(context, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
