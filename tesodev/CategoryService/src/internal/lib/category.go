package lib

import (
	"time"

	"github.com/turgut-nergin/tesodev_work1/internal/models"
)

func ResponseAssign(category *models.Category) *models.CategoryResponse {
	return &models.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,

		CreatedAt: time.Unix(category.CreatedAt, 0),
		UpdatedAt: time.Unix(category.UpdatedAt, 0),
	}

}
