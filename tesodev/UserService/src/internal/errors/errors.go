package errors

import "github.com/labstack/echo/v4"

type Error struct {
	ApplicationName string `json:"applicationName"`
	Operation       string `json:"operation"`
	Description     string `json:"description"`
	StatusCode      int    `json:"statusCode"`
	ErrorCode       int    `json:"errorCode"`
}

func New(applicationName, operation, description string, statusCode, errorCode int) *Error {
	return &Error{
		ApplicationName: applicationName,
		Operation:       operation,
		Description:     description,
		StatusCode:      statusCode,
		ErrorCode:       errorCode,
	}
}

func (e *Error) WrapOperation(operation string) *Error {
	e.Operation = operation
	return e
}

func (e *Error) ToResponse(ctx echo.Context) error {
	return ctx.JSON(e.StatusCode, e)
}

func (e *Error) WrapDesc(desc string) *Error {
	e.Description = desc
	return e
}

func (e *Error) WrapErrorCode(errorCode int) *Error {
	e.ErrorCode = errorCode
	return e
}

var (
	UnknownError    = New("user", "", "unknown error", 500, 0)
	ValidationError = New("user", "validation", "validation error", 400, 0)
	NotFound        = New("user", "handler", "not found user", 404, 0)
	FindFailed      = New("category", "Repository", "Categories find failed! ", 500, 0)
)
