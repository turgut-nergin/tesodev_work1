package models

import (
	"time"
)

type UserResponse struct {
	UserId    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type User struct {
	UserId    string `bson:"userId,omitempty"`
	UserName  string `bson:"userName"`
	Password  string `bson:"password"`
	Email     string `bson:"email"`
	Type      string `bson:"type"`
	CreatedAt int64  `bson:",omitempty"`
	UpdatedAt int64  `bson:",omitempty"`
}

type UserRequest struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Type     string `json:"type,omitempty"`
}
