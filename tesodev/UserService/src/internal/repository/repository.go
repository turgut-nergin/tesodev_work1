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

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result, err := r.db.DeleteOne(context, bson.M{"_id": id})
	return result.DeletedCount, err
}

func (r *Repository) Upsert(id string, user *models.User) *models.UpSertResult {
	filter := bson.M{"_id": id}
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
	upsertResult.ID = id
	upsertResult.ModifiedCount = result.ModifiedCount

	return &upsertResult
}

func (r *Repository) Find(limit, offset int64, filter map[string]interface{}, sortField string, sortDirection int) (*models.UserRows, *errors.Error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	countChan := make(chan models.CountOrError)

	go func() {
		totalCount, err := r.db.CountDocuments(context, filter)
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

	options = options.SetSkip(offset * limit).SetLimit(limit) //pagination set

	cur, err := r.db.Find(context, filter, options)

	if err != nil {
		return nil, errors.FindFailed.WrapErrorCode(4001)
	}

	var users []*models.User

	err = cur.All(context, &users)

	if err != nil {
		return nil, errors.UnknownError.WrapErrorCode(4002)
	}

	var userResponse []models.UserResponse

	for _, user := range users {
		ticketResponse := lib.ResponseAssign(user)
		userResponse = append(userResponse, *ticketResponse)
	}

	countOrError := <-countChan

	if countOrError.Error != nil {
		return nil, countOrError.Error
	}

	return &models.UserRows{
		RowCount: countOrError.TotalCount,
		Users:    userResponse,
	}, nil
}

func (r *Repository) FindByUserName(userName string) (*models.User, error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"userName": userName,
	}
	user := models.User{}

	err := r.db.FindOne(context, filter).Decode(&user)

	return &user, err
}
