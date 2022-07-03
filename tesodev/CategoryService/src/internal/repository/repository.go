package repository

import (
	"context"
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
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
	category := models.Category{}

	if err := r.collection.FindOne(context, bson.M{"_id": id}).Decode(&category); err != nil {
		return nil, err
	}

	return &category, nil
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

func (r *Repository) Update(category *models.Category) (*int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": category.Id}
	result, err := r.collection.UpdateOne(context, filter, category)
	if err != nil {
		return nil, err
	}

	return &result.ModifiedCount, nil

}

// func (r *Repository) Find(limit int64, offset int64, filter []primitive.D) (*int64, error) {
// 	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	filter := bson.M{"_id": category.Id}
// 	result, err := r.collection.UpdateOne(context, filter, category)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &result.ModifiedCount, nil

// }
