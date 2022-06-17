package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	db *mongo.Collection
}

func New(mongoClient *mongo.Collection) *Repository {
	driverRepository := Repository{mongoClient}
	return &driverRepository
}

func (r *Repository) Get(id string) (*user.User, error) {
	user := user.User{}

	if err := r.db.FindOne(context.Background(), bson.M{"id": id}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Delete(id string) (int64, error) {
	result, err := r.db.DeleteOne(context.Background(), bson.M{"id": id})
	return result.DeletedCount, err
}

func (r *Repository) Upsert(id string, user *user.User) (*int64, *string, error) {
	filter := bson.M{"id": id}
	userId := uuid.New().String()
	createdAt := lib.TimeStampNow()
	update := bson.M{
		"$setOnInsert": bson.M{
			"createdAt": createdAt,
			"id":        userId,
		},
		"$set": user}
	opts := options.Update().SetUpsert(true)
	result, err := r.db.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, nil, err
	}

	return &result.ModifiedCount, &userId, nil
}
