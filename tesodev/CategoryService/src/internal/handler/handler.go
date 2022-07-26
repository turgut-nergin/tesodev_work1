package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/turgut-nergin/tesodev_work1/config"
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
	cfg        *config.Config
}

func New(repository *repository.Repository, config *config.Config) *Handler {
	return &Handler{repository: repository, cfg: config}
}

// GetCategory
// @Summary  Get Category by Id
// @Description Get Category by ID
// @Tags cateroies
// @Accept json
// @Param Token header string true "Bearer"
// @Produce json
// @Param categoryId query string false "categoryId"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.CategoryResponse
// @Router /user/category [GET]
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
// @Param Token header string true "Bearer"
// @Accept json
// @Produce json
// @Param models.CategoryRequest body models.CategoryRequest true "For Create a Categry"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /admin/category [POST]
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
		return errors.UnknownError.WrapErrorCode(4014).WrapDesc(err.Error()).ToResponse(c)
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
// @Param categoryId path string true "Category Id"
// @Param Token header string true "Bearer"
// @Param models.CategoryRequest body models.CategoryRequest true "For Update a Categry"
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /user/category/{categoryId} [PUT]
func (h *Handler) UpdateCategory(c echo.Context) error {
	id := c.Param("categoryId")

	if _, err := uuid.Parse(id); err != nil {
		return errors.ValidationError.WrapErrorCode(1008).WrapDesc(err.Error()).ToResponse(c)
	}

	categoryReq := models.CategoryRequest{}

	if err := json.NewDecoder(c.Request().Body).Decode(&categoryReq); err != nil {
		return errors.ValidationError.WrapErrorCode(1009).WrapDesc(err.Error()).ToResponse(c)
	}

	category := models.Category{}

	category.Id = id
	category.UpdatedAt = lib.TimeStampNow()
	category.Name = categoryReq.Name

	modifiedCount, error := h.repository.Update(&category)

	if error != nil {
		return error.ToResponse(c)
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
// @Param categoryId path string true "Category Id"
// @Param Token header string true "Bearer"
// @Failure 404 {object} bool
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} bool
// @Router /admin/category/{categoryId} [delete]
func (h *Handler) DeleteCategory(c echo.Context) error {
	id := c.Param("categoryId")
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

	return c.JSON(http.StatusOK, true)
}

// GetCategories
// @Summary  Get Categories by params
// @Description Get Categories by params
// @Tags cateroies
// @Accept json
// @Produce json
// @Param name query string false "name"
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param sort query string false "sort"
// @Param direction query string false "direction"
// @Param Token header string true "Bearer"
// @Failure 404 {object} errors.Error
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Succes 200 {object} models.CategoryResponse
// @Router /user/categories [get]
func (h *Handler) GetCategories(c echo.Context) error {

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, offset := lib.ValidatePaginator(limitStr, offsetStr, h.cfg.MaxPageLimit)
	filter := map[string]interface{}{}

	if name := c.QueryParam("name"); name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{
			Pattern: name,
			Options: "i",
		}}
	}

	sortField := lib.GetAcceptedSortField(c.QueryParam("sort"))              //for example name
	sortDirection := lib.GetAcceptedSortDirection(c.QueryParam("direction")) //asc desc -> 0,-1
	categories, err := h.repository.Find(limit, offset, filter, sortField, sortDirection)

	if err != nil {
		response := models.CategoryRows{
			RowCount:   0,
			Categories: nil,
		}

		return c.JSON(200, response)
	}

	return c.JSON(200, categories)
}

// CSV
// @Summary Read csv file
// @Description Read csv file
// @Tags CSV
// @Accept json
// @Param file formData file true  "this is a test file"
// @Produce json
// @Failure 404 {object} errors.Error
// @Succes 200 {object} models.User
// @Router /csv [post]
func (h *Handler) ReadCsv(c echo.Context) error {
	concurrentCounter := 10
	now := time.Now() // current local time
	// 1 second
	workQueue := make(chan string)

	complete := make(chan []string)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()

	if err != nil {
		return err
	}

	defer file.Close()

	go readFile(file, workQueue)

	// Now read them all off, concurrently.
	for i := 0; i < concurrentCounter; i++ {
		go startWorking(workQueue, complete)
	}

	var result []string
	// Wait for everyone to finish.
	for i := 0; i < concurrentCounter; i++ {
		result = append(result, <-complete...)
	}

	// fileHeader, err := c.FormFile("file")
	// if err != nil {
	// 	return err
	// }

	// file, err := fileHeader.Open()

	// if err != nil {
	// 	return err
	// }

	// defer file.Close()
	// data, err := csv.NewReader(file).ReadAll()
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, err)
	// }

	// for _, line := range data {
	// 	fmt.Println(line)

	// }
	elapsed := time.Since(now)
	fmt.Println(result)
	return c.JSON(http.StatusOK, elapsed)

}

func startWorking(queue chan string, complete chan []string) {
	var result []string
	for line := range queue {
		result = append(result, line)
	}

	// Let the main process know we're done.
	complete <- result
}

func readFile(file multipart.File, workQueue chan<- string) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		workQueue <- scanner.Text() //read line by line sending channel
	}

	// Close the channel so everyone reading from it knows we're done.
	close(workQueue)

}
