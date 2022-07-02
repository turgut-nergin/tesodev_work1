package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
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

func (r *Repository) Get(id string) (*models.User, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	user := models.User{}

	if err := r.db.FindOne(context, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Delete(id string) (int64, error) {
	result, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return result.DeletedCount, err
}

func (r *Repository) Upsert(id string, user *models.User) *models.UpSertResult {
	filter := bson.M{"_id": id}
	userId := uuid.New().String()
	createdAt := lib.TimeStampNow()

	update := bson.M{
		"$setOnInsert": bson.M{
			"createdAt": createdAt,
		},
		"$set": user}

	opts := options.Update().SetUpsert(true)
	result, err := r.db.UpdateOne(context.TODO(), filter, update, opts)

	upsertResult := models.UpSertResult{}

	if err != nil {
		upsertResult.Err = err
		upsertResult.ErrCode = 1010
		return &upsertResult
	}

	upsertResult.Err = err
	upsertResult.ID = userId
	upsertResult.ModifiedCount = result.ModifiedCount

	return &upsertResult
}
