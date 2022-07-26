package repository

import (
	"context"

	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"

	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	collection *mongo.Collection
}

func New(mongoClient *mongo.Collection) *Repository {
	driverRepository := Repository{mongoClient}
	return &driverRepository
}

func (r *Repository) FindOne(id string) (*models.Category, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	var category *models.Category

	if err := r.collection.FindOne(context, bson.M{"_id": id}).Decode(&category); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *Repository) Delete(id string) (int64, error) {
	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return result.DeletedCount, err
}

func (r *Repository) Insert(category *models.Category) (*string, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := r.collection.InsertOne(context, category); err != nil {
		return nil, err
	}

	return &category.Id, nil

}

func (r *Repository) Update(category *models.Category) (*int64, *errors.Error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": category.Id}
	result, err := r.collection.UpdateOne(context, filter, bson.M{"$set": category})
	if err != nil {
		return nil, errors.UnknownError.WrapOperation("repository").WrapErrorCode(4000).WrapDesc(err.Error())
	}

	return &result.ModifiedCount, nil

}

func (r *Repository) Find(limit, offset int64, filter map[string]interface{}, sortField string, sortDirection int) (*models.CategoryRows, *errors.Error) {

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

	options = options.SetSkip(limit * offset).SetLimit(limit) //pagination set

	cur, err := r.collection.Find(context, filter, options)

	if err != nil {
		return nil, errors.FindFailed.WrapErrorCode(4001)
	}

	var categories []models.Category

	err = cur.All(context, &categories)

	if err != nil {
		return nil, errors.UnknownError.WrapErrorCode(4002)
	}

	countOrError := <-countChan

	if countOrError.Error != nil {
		return nil, countOrError.Error
	}

	return &models.CategoryRows{
		RowCount:   countOrError.TotalCount,
		Categories: &categories,
	}, nil
}
