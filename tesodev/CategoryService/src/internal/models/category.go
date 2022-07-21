package models

import (
	"time"
)

type CategoryResponse struct {
	Id        string    `json:"_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type CategoryRows struct {
	RowCount   int64       `json:"rowCount"`
	Categories *[]Category `json:"categories"`
}

type Category struct {
	Id        string `bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name,omitempty"`
	CreatedAt int64  `bson:",omitempty"`
	UpdatedAt int64  `bson:",omitempty"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}
