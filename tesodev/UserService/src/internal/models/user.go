package models

import (
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
)

type UserResponse struct {
	UserId    string    `json:"_id"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type User struct {
	UserId    string `bson:"_id,omitempty"`
	UserName  string `bson:"userName,omitempty"`
	Password  string `bson:"password,omitempty"`
	Email     string `bson:"email,omitempty"`
	Type      string `bson:"type,omitempty"`
	CreatedAt int64  `bson:",omitempty"`
	UpdatedAt int64  `bson:",omitempty"`
}

type UserRequest struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Type     string `json:"type,omitempty"`
}

type UserRows struct {
	RowCount int64          `json:"rowCount"`
	Users    []UserResponse `json:"users"`
}

type Authentication struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	UserName    string `json:"userName"`
	TokenString string `json:"token"`
}

type CountOrError struct {
	TotalCount int64
	Error      *errors.Error
}
