package models

import "time"

type Answer struct {
	Id        string `bson:"_id,omitempty"`
	Body      string `bson:"body" json:"body"`
	CreatedBy string `bson:"createdBy,omitempty"`
	CreatedAt int64  `bson:"createdAt,omitempty"`
	UpdatedAt int64  `bson:"updatedAt,omitempty"`
}

type AnswerResponse struct {
	Id        string    `json:"_id,omitempty"`
	Body      string    `json:"ticketId"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
