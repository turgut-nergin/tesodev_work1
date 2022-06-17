package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models/request"
	"github.com/turgut-nergin/tesodev_work1/internal/models/user"
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
	user, err := h.repository.Get(id)
	if err != nil {
		return errors.NotFound.WrapErrorCode(1000).
			WrapDesc(fmt.Sprintf("User id: %v not found", id)).ToResponse(c)

	}

	respUser := lib.ResponseAssign(user)

	return c.JSON(http.StatusOK, respUser)
}

func (h *Handler) UpsertUser(c echo.Context) error {
	id := c.QueryParam("userId")

	userReq := request.User{}

	json.NewDecoder(c.Request().Body).Decode(&userReq)
	err := userReq.Validate()

	if err != nil {
		return errors.ValidationError.WrapErrorCode(1003).WrapDesc(err.Error()).ToResponse(c)
	}

	user := user.User{
		UserName:  userReq.UserName,
		Password:  userReq.Password,
		Email:     userReq.Email,
		Type:      userReq.Type,
		UpdatedAt: lib.TimeStampNow(),
	}

	modifiedCount, newId, err := h.repository.Upsert(id, &user)
	if err != nil {
		return errors.UnknownError.WrapErrorCode(1005).WrapDesc(err.Error()).ToResponse(c)
	}
	if *modifiedCount == 1 {
		return c.JSON(http.StatusOK, id)
	}

	return c.JSON(http.StatusCreated, newId)

}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("userId")
	deleteResult, err := h.repository.Delete(id)
	if err != nil {
		errors.UnknownError.WrapErrorCode(1000).
			WrapDesc(err.Error()).ToResponse(c)
	}
	if deleteResult == 0 {
		return c.JSON(http.StatusNotFound, false)
	}
	return c.JSON(http.StatusOK, true)
}
