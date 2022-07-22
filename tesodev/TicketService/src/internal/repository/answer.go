package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type AnswerRepository interface {
	GetAnswers(ticketId string) ([]models.Answer, *errors.Error)
	GetAnswer(answerId string) (*models.Answer, *errors.Error)
	DeleteAnswer(ticketId string) (int64, error)
	UpdateAnswer(answer *models.Answer) (*int64, error)
	CreateAnswer(answer *models.Answer) (*string, error)
}

func NewAnswer(mongoClient *mongo.Collection) *Repository {
	return &Repository{mongoClient}
}

func (r *Repository) GetAnswers(ticketId string) ([]models.Answer, *errors.Error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	cursor, err := r.collection.Find(context, bson.M{"ticketId": ticketId})

	if err != nil {
		return nil, errors.UnknownError.WrapOperation("repository").WrapErrorCode(4060).WrapDesc(err.Error())
	}

	var answers []models.Answer
	if err = cursor.All(context, &answers); err != nil {
		return nil, errors.UnknownError.WrapOperation("repository").WrapErrorCode(4061).WrapDesc(err.Error())
	}
	return answers, nil
}

func (r *Repository) GetAnswer(answerId string) (*models.Answer, *errors.Error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	answer := models.Answer{}

	if err := r.collection.FindOne(context, bson.M{"_id": answerId}).Decode(&answer); err != nil {
		return nil, errors.UnknownError.WrapOperation("repository").WrapErrorCode(4060).WrapDesc(err.Error())
	}

	fmt.Println("............................")
	fmt.Println(answer)
	return &answer, nil
}

func (r *Repository) UpdateAnswer(answer *models.Answer) (*int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	filter := bson.M{"_id": answer.Id}

	result, err := r.collection.UpdateOne(context, filter, bson.M{
		"$set": answer})
	if err != nil {
		return nil, err
	}

	return &result.ModifiedCount, nil
}

func (r *Repository) CreateAnswer(answer *models.Answer) (*string, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := r.collection.InsertOne(context, answer); err != nil {
		return nil, err
	}

	return &answer.Id, nil

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
