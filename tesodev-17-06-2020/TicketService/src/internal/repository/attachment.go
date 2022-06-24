package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type AttachmentRepository interface {
	InsertAttachment(tickedId string, attachments []models.Attachment) error
	DeleteAttachment(ticketId string) (int64, error)
	GetAttachment(ticketId string) ([]models.Attachment, error)
}

func NewAttachment(mongoClient *mongo.Database) *Repository {
	return &Repository{mongoClient.Collection("Attachments")}
}

func (r *Repository) InsertAttachment(tickedId string, attachments []models.Attachment) error {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	var attachmentList []interface{}

	for _, attachment := range attachments {
		attachment.Id = uuid.New().String()
		attachment.TicketId = tickedId
		attachmentList = append(attachmentList, attachment)
	}

	if _, err := r.collection.InsertMany(context, attachmentList); err != nil {
		return nil
	}

	return nil
}

func (r *Repository) DeleteAttachment(ticketId string) (int64, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	result, err := r.collection.DeleteMany(context, bson.M{"ticketId": ticketId})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *Repository) GetAttachment(ticketId string) ([]models.Attachment, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	attachments := []models.Attachment{}

	cursor, err := r.collection.Find(context, bson.M{"ticketId": ticketId})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context, &attachments); err != nil {
		return nil, err
	}

	return attachments, nil
}
