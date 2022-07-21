package models

import (
	"time"
)

type TicketResponse struct {
	Id             string           `json:"_id"`
	Subject        string           `json:"subject"`
	Body           string           `json:"body"`
	Category       *Category        `json:"category"`
	Attachments    []Attachment     `json:"attachments"`
	Answers        []AnswerResponse `json:"answers,omitempty"`
	CreatedBy      string           `json:"createdBy"`
	Status         string           `json:"status"`
	LastAnsweredAt time.Time        `json:"lastAnsweredAt,omitempty"`
	CreatedAt      time.Time        `json:"createdAt,omitempty"`
	UpdatedAt      time.Time        `json:"updatedAt,omitempty"`
}

type TicketRows struct {
	RowCount int64            `json:"rowCount"`
	Tickets  []TicketResponse `json:"tickets"`
}

type Ticket struct {
	Id             string       `bson:"_id,omitempty"`
	Subject        string       `bson:"subject,omitempty"`
	Body           string       `bson:"body,omitempty"`
	CategoryId     string       `bson:"categoryId,omitempty"`
	CreatedBy      string       `bson:"createdBy,omitempty"`
	Attachments    []Attachment `bson:"attachment,omitempty"`
	Status         string       `bson:"status,omitempty"`
	LastAnsweredAt int64        `bson:"lastAnsweredAt,omitempty"`
	CreatedAt      int64        `bson:"createdAt,omitempty"`
	UpdatedAt      int64        `bson:"updatedAt,omitempty"`
}

type TicketRequest struct {
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Attachments []Attachment `json:"attachments,inline"`
	Status      string       `json:"status"`
}

type Attachment struct {
	FileName string `bson:"fileName" json:"fileName"`
	FilePath string `bson:"filePath" json:"filePath"`
}

type Category struct {
	Id        string    `json:"_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Answer struct {
	Id        string `bson:"_id,omitempty"`
	Body      string `bson:"body" json:"body"`
	TicketId  string `bson:"ticketId,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	CreatedAt int64  `bson:"createdAt,omitempty"`
	UpdatedAt int64  `bson:"updatedAt,omitempty"`
}

type AnswerResponse struct {
	Id        string    `json:"_id,omitempty"`
	Body      string    `json:"body"`
	UserId    string    `json:"userId"`
	TicketId  string    `json:"ticketId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
