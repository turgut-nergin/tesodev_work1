package repository

import (
	"context"
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type TicketRepository interface {
	InsertTicket(ticket *models.Ticket) (*string, *errors.Error)
	DeleteTicket(id string) (int64, error)
	GetTicket(id string) (*models.Ticket, error)
	Update(ticket models.Ticket) (*int64, *errors.Error)
	Find(limit int64, offset int64, filter map[string]interface{}, sortField string, sortDirection int) (*models.TicketRows, *errors.Error)
}

func NewTicket(mongoClient *mongo.Collection) *Repository {
	return &Repository{mongoClient}
}

func (r *Repository) InsertTicket(ticket *models.Ticket) (*string, *errors.Error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if _, err := r.collection.InsertOne(context, ticket); err != nil {

		return nil, errors.UnknownError.WrapOperation("repository").WrapErrorCode(4062).WrapDesc(err.Error())
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

func (r *Repository) Update(ticket models.Ticket) (*int64, *errors.Error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	filter := bson.M{"_id": ticket.Id}

	result, err := r.collection.UpdateOne(context, filter, bson.M{"$set": ticket})
	if err != nil {
		return nil, errors.UnknownError.WrapErrorCode(4062).WrapDesc(err.Error()).WrapOperation("repository")
	}
	return &result.ModifiedCount, nil
}

func (r *Repository) Find(limit, offset int64, filter map[string]interface{}, sortField string, sortDirection int) (*models.TicketRows, *errors.Error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	countChan := make(chan models.CountOrError)

	go func() {
		totalCount, err := r.collection.CountDocuments(context, filter)
		countOrError := models.CountOrError{
			TotalCount: totalCount,
			Error:      nil,
		}
		if err != nil {
			countOrError.Error = errors.FindFailed.WrapErrorCode(4000)
		}

		countChan <- countOrError
	}()

	options := options.Find()

	if sortField != "" && sortDirection == 0 {
		options = options.SetSort(bson.D{{sortField, sortDirection}})
	}

	options = options.SetLimit(limit).SetSkip(offset * limit) //pagination set

	cur, err := r.collection.Find(context, filter, options)

	if err != nil {
		return nil, errors.FindFailed.WrapErrorCode(4001)
	}

	var tickets []*models.Ticket

	err = cur.All(context, &tickets)

	if err != nil {
		return nil, errors.UnknownError.WrapErrorCode(4002)
	}

	var ticketsResponse []models.TicketResponse

	for _, ticket := range tickets {
		ticketResponse := lib.TicketResponseAssign(ticket)
		ticketsResponse = append(ticketsResponse, *ticketResponse)
	}

	countOrError := <-countChan

	if countOrError.Error != nil {
		return nil, countOrError.Error
	}

	return &models.TicketRows{
		RowCount: countOrError.TotalCount,
		Tickets:  ticketsResponse,
	}, nil
}
