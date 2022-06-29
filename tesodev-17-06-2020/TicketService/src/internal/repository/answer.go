package repository

import (
	"context"
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type AnswerRepository interface {
	GetAnswers(ticketId string) ([]models.Answer, error)
	UpsertAnswer(id string, answer *models.Answer) (*int64, error)
	DeleteAnswer(ticketId string) (int64, error)
}

func NewAnswer(mongoClient *mongo.Collection) *Repository {
	return &Repository{mongoClient}
}

func (r *Repository) GetAnswers(ticketId string) ([]models.Answer, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	cursor, err := r.collection.Find(context, bson.M{"createdBy": ticketId})

	if err != nil {
		return nil, err
	}

	var answers []models.Answer
	if err = cursor.All(context, &answers); err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *Repository) UpsertAnswer(id string, answer *models.Answer) (*int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	filter := bson.M{"_id": id}

	update := bson.M{
		"$setOnInsert": bson.M{
			"createdAt": lib.TimeStampNow(),
		},
		"$set": answer}
	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(context, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return &result.ModifiedCount, nil
}

func (r *Repository) DeleteAnswer(ticketId string) (int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	result, err := r.collection.DeleteMany(context, bson.M{"createdBy": ticketId})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
