package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repository *repository.Repository
	cfg        *config.Config
}

func New(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

// GetCategory
// @Summary  Get Category by Id
// @Description Get Category by ID
// @Tags cateroies
// @Accept json
// @Produce json
// @Param categoryId query string false "categoryId"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.CategoryResponse
// @Router /category [GET]
func (h *Handler) GetCategory(c echo.Context) error {

	id := c.QueryParam("categoryId")

	if _, err := uuid.Parse(id); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	category, err := h.repository.FindOne(id)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(1000).
			WrapDesc(fmt.Sprintf(err.Error())).ToResponse(c)

	}
	if category == nil {
		return errors.NotFound.WrapErrorCode(11000).
			WrapDesc(fmt.Sprintf("Category id: %v not found", id)).ToResponse(c)

	}
	responseCategory := lib.ResponseAssign(category)

	return c.JSON(http.StatusOK, responseCategory)
}

// CreateCategory
// @Summary  Create Category
// @Description Create Category
// @Tags cateroies
// @Accept json
// @Produce json
// @Param models.CategoryRequest body models.CategoryRequest true "For Create a Categry"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /category [POST]
func (h *Handler) CreateCategory(c echo.Context) error {

	categoryReq := models.CategoryRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&categoryReq); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	category := models.Category{}

	category.Id = uuid.New().String()
	category.CreatedAt = lib.TimeStampNow()
	category.Name = categoryReq.Name
	categoryId, err := h.repository.Insert(&category)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(4000).WrapDesc(err.Error()).ToResponse(c)
	}

	return c.JSON(http.StatusCreated, categoryId)

}

// UpdateCategory by ID
// @Summary  Update Category
// @Description Update Category by ID
// @Tags cateroies
// @Accept json
// @Produce json
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /category/{categoryId} [PUT]
func (h *Handler) UpdateCategory(c echo.Context) error {
	id := c.Param("categoryId")

	if _, err := uuid.Parse(id); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	category := models.Category{}

	if err := json.NewDecoder(c.Request().Body).Decode(&category); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}
	category.Id = id
	category.UpdatedAt = lib.TimeStampNow()
	modifiedCount, err := h.repository.Update(&category)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(4000).WrapDesc(err.Error()).ToResponse(c)
	}

	if *modifiedCount == 0 {
		return c.JSON(http.StatusNotFound, false)

	}

	return c.JSON(http.StatusOK, true)

}

// DeleteCategory by ID
// @Summary  Delete Category
// @Description Delete Category by ID
// @Tags cateroies
// @Accept json
// @Produce json
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /category/{categoryId} [delete]
func (h *Handler) DeleteCategory(c echo.Context) error {
	id := c.Param("categoryId")
	_, err := uuid.Parse(id)

	if err == nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	deleteResult, err := h.repository.Delete(id)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(1000).
			WrapDesc(err.Error()).ToResponse(c)
	}

	if deleteResult == 0 {
		return c.JSON(http.StatusNotFound, false)
	}

	return c.JSON(http.StatusOK, true)
}

// func (h *Handler) GetCategories(c echo.Context) error {

// 	limitStr := c.QueryParam("limit")
// 	offsetStr := c.QueryParam("offset")
// 	offset, limit := lib.ValidatePaginator(limitStr, offsetStr, int(h.cfg.MaxPageLimit))
// 	filter := []bson.D{}

// 	if name := c.QueryParam("name"); name != "" && len(name) < 10 {
// 		filter = append(filter, bson.D{{"$regex", primitive.Regex{
// 			Pattern: name,
// 			Options: "i",
// 		}}})
// 	}
// 	categories, err := h.repository.Find(limit, offset, filter)

// 	if err != nil {

// 	}

// }
