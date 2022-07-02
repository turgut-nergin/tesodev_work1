package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"

	"github.com/labstack/echo"
)

type Handler struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetUser(c echo.Context) error {

	id := c.Param("userId")

	if _, err := uuid.Parse(id); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	user, err := h.repository.Get(id)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(1000).
			WrapDesc(fmt.Sprintf(err.Error())).ToResponse(c)

	}
	if user == nil {
		return errors.NotFound.WrapErrorCode(11000).
			WrapDesc(fmt.Sprintf("User id: %v not found", id)).ToResponse(c)

	}
	respUser := lib.ResponseAssign(user)

	return c.JSON(http.StatusOK, respUser)
}

func (h *Handler) UpsertUser(c echo.Context) error {
	id := c.QueryParam("userId")

	if id != "" {
		if _, err := uuid.Parse(id); err != nil {
			return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
		}
	} else {
		id = uuid.New().String()
	}

	userReq := models.UserRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&userReq); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	if id == "" {
		if err := userReq.ValidateInsert(); err != nil {
			return errors.ValidationError.WrapErrorCode(1003).WrapDesc(err.Error()).ToResponse(c)
		}
	} else {
		if err := userReq.ValidateUpdate(); err != nil {
			return errors.ValidationError.WrapErrorCode(1004).WrapDesc(err.Error()).ToResponse(c)
		}
	}

	hashedPass, err := lib.HashPassword(userReq.Password)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(1006).WrapDesc(err.Error()).ToResponse(c)
	}

	user := models.User{
		UserName:  userReq.UserName,
		Password:  hashedPass,
		Email:     userReq.Email,
		Type:      userReq.Type,
		UpdatedAt: lib.TimeStampNow(),
	}

	result := h.repository.Upsert(id, &user)
	if result.Err != nil {
		return errors.UnknownError.WrapErrorCode(result.ErrCode).WrapDesc(result.Err.Error()).ToResponse(c)
	}

	if result.ModifiedCount == 1 {
		return c.JSON(http.StatusOK, id)
	}

	return c.JSON(http.StatusCreated, result.ID)

}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("userId")
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
