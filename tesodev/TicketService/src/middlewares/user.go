package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/turgut-nergin/tesodev_work1/internal/errors"
)

func IsAuthorized(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header["Token"]
			if tokenString == nil {
				return errors.NotFound.WrapOperation("middleware").WrapErrorCode(1008).WrapDesc("No Token Found").ToResponse(c)
			}

			var apiKey = []byte("tesodev")

			token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error in parsing")
				}
				return apiKey, nil

			})

			if err != nil {
				return errors.ValidationError.WrapOperation("middleware").WrapErrorCode(5001).WrapDesc("Your Token has been expired").ToResponse(c)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if claims["role"] == role {
					c.Response().Header().Set("role", "admin")
					return next(c)
				}
			}
			return c.JSON(http.StatusUnauthorized, "unAuthorized")

		}
	}
}
