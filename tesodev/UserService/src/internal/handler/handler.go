package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/config"
	"github.com/turgut-nergin/tesodev_work1/internal/broker"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/lib"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/turgut-nergin/tesodev_work1/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repository *repository.Repository
	cfg        *config.MongoConfig
	rabbitmq   *models.RabbitMQ
}

func New(repository *repository.Repository, config *config.MongoConfig, rabbitMQModel models.RabbitMQ) *Handler {
	return &Handler{
		repository: repository,
		cfg:        config,
		rabbitmq:   &rabbitMQModel,
	}
}

// GetUser
// @Summary get user by id
// @Description get id
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "user ID"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.UserResponse
// @Router /user/{userId} [get]
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

// UpsertUser
// @Summary Update and Create User
// @Description Create User and Update User
// @Tags users
// @Accept json
// @Produce json
// @Param userId query string false "user ID"
// @Param models.UserRequest body models.UserRequest true "For upsert an User"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 201 {object} models.UpSertResult
// @Succes 200 {object} string
// @Router /user [post]
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

// DeleteUser by ID
// @Summary  Delete User
// @Description Delete User
// @Tags users
// @Accept json
// @Produce json
// @Failure 404 {object} bool
// @Param userId path string true "User Id"
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /user/{userId} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("userId")
	_, err := uuid.Parse(id)

	if err != nil {
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

	h.rabbitmq.Message = id
	broker.PublishConsumer(h.rabbitmq)

	return c.JSON(http.StatusOK, true)
}

// DeleteUser by ID
// @Summary  Validate User
// @Description Validate User
// @Tags users
// @Accept json
// @Produce json
// @Param userId query string false "userId"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /user [get]
func (h *Handler) Validate(c echo.Context) error {
	id := c.QueryParam("userId")
	_, err := uuid.Parse(id)

	if err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	userInfo, err := h.repository.Get(id)

	if err != nil {
		return errors.UnknownError.WrapErrorCode(1000).
			WrapDesc(err.Error()).ToResponse(c)
	}

	if userInfo == nil {
		return c.JSON(http.StatusNotFound, false)
	}

	return c.JSON(http.StatusOK, true)
}

// GetUsers
// @Summary  Get Users by params
// @Description Get Users by params
// @Tags Users
// @Accept json
// @Produce json
// @Param userName query string false "userName"
// @Param email query string false "email"
// @Param type query string false "type"
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param sort query string false "sort"
// @Param direction query string false "direction"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.CategoryResponse
// @Router /users [get]
func (h *Handler) GetUsers(c echo.Context) error {

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, offset := lib.ValidatePaginator(limitStr, offsetStr, h.cfg.MaxPageLimit)

	filter := map[string]interface{}{}

	if userName := c.QueryParam("userName"); userName != "" {
		filter["userName"] = bson.M{"$regex": primitive.Regex{
			Pattern: userName,
			Options: "i",
		}}
	}

	if email := c.QueryParam("email"); email != "" {
		filter["email"] = bson.M{"$regex": primitive.Regex{
			Pattern: email,
			Options: "i",
		}}
	}

	if userType := c.QueryParam("type"); userType != "" {
		filter["type"] = bson.M{"$regex": primitive.Regex{
			Pattern: userType,
			Options: "i",
		}}
	}

	sortField := lib.GetAcceptedSortField(c.QueryParam("sort"))              //for example name
	sortDirection := lib.GetAcceptedSortDirection(c.QueryParam("direction")) //asc desc -> 0,-1
	tickets, err := h.repository.Find(limit, offset, filter, sortField, sortDirection)
	if err != nil {
		return err.ToResponse(c)
	}

	return c.JSON(http.StatusOK, tickets)
}

// Login
// @Summary  Login by username and password
// @Description Login by username and password
// @Tags Login
// @Accept json
// @Produce json
// @Param models.UserRequest body models.Authentication true "To Login"
// @Failure 404 {object} errors.Error
// @Succes 200 {object} models.User
// @Router /login [post]
func (h *Handler) Login(c echo.Context) error {

	login := models.Authentication{}

	if err := json.NewDecoder(c.Request().Body).Decode(&login); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	result, err := h.repository.FindByUserName(login.UserName)

	if err != nil {
		return c.JSON(http.StatusNotFound, result)
	}

	Check := lib.DoPasswordsMatch(result.Password, login.Password)

	if !Check {
		return c.JSON(http.StatusUnauthorized, "Authentication error")
	}

	validToken, err := lib.GenerateJWT(result.UserName, result.Type)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Failed to generate token")

	}
	var token models.Token
	token.UserName = result.UserName
	token.Role = result.Type
	token.TokenString = validToken
	return c.JSON(http.StatusOK, token)

}

// isAuthorized
// @Summary  isAuthorized by token
// @Description isAuthorized by token
// @Tags isAuthorized
// @Accept json
// @Produce json
// @Failure 404 {object} errors.Error
// @Succes 200 {object} models.User
// @Param Token header string true "Bearer"
// @Router /isAuthorized [get]
func (h *Handler) IsAuthorized(c echo.Context) error {

	tokenString := c.Request().Header["Token"]
	if tokenString == nil {
		return errors.NotFound.WrapErrorCode(1008).WrapDesc("No Token Found").ToResponse(c)
	}

	var apiKey = []byte("tesodev")

	token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return apiKey, nil

	})

	if err != nil {
		return errors.ValidationError.WrapErrorCode(5001).WrapDesc("Your Token has been expired").ToResponse(c)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if claims["role"] == "admin" {

			c.Response().Header().Set("role", "admin")
			return c.JSON(http.StatusOK, "admin")

		} else if claims["role"] == "user" {
			c.Response().Header().Set("role", "user")
			return c.JSON(http.StatusOK, "user")

		}
	}
	return c.JSON(http.StatusUnauthorized, "unAuthorized")

}
