package lib

import (
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/models/response"
	"github.com/turgut-nergin/tesodev_work1/internal/models/user"
)

func ResponseAssign(user *user.User) *response.User {
	return &response.User{
		ID:        user.ID,
		UserName:  user.UserName,
		Password:  user.Password,
		Email:     user.Email,
		Type:      user.Type,
		CreatedAt: time.Unix(user.CreatedAt, 0),
		UpdatedAt: time.Unix(user.UpdatedAt, 0),
	}

}
