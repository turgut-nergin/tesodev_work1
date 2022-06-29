package models

import (
	"time"
)

type TicketResponse struct {
	Id             string       `json:"_id"`
	Subject        string       `json:"subject"`
	Body           string       `json:"body"`
	Category       Category     `json:"category"`
	Attachments    []Attachment `json:"attachments"`
	Answers        []Answer     `json:"answers,omitempty"`
	CreatedBy      string       `json:"createdBy"`
	Status         string       `json:"status"`
	LastAnsweredAt time.Time    `json:"lastAnsweredAt,omitempty"`
	CreatedAt      time.Time    `json:"createdAt,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt,omitempty"`
}

type Ticket struct {
	Id         string `bson:"_id,omitempty"`
	Subject    string `bson:"subject,omitempty"`
	Body       string `bson:"body,omitempty"`
	CategoryId string `bson:"categoryId,omitempty"`
	CreatedBy  string `bson:"createdBy,omitempty"`

	Status         string `bson:"status,omitempty"`
	LastAnsweredAt int64  `bson:"lastAnsweredAt,omitempty"`
	CreatedAt      int64  `bson:"createdAt,omitempty"`
	UpdatedAt      int64  `bson:"updatedAt,omitempty"`
}

type TicketRequest struct {
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Attachments []Attachment `json:"attachments,inline"`
	Status      string       `json:"status"`
}
