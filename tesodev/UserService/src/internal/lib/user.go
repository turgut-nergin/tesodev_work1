package lib

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func ResponseAssign(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		UserId:    user.UserId,
		UserName:  user.UserName,
		Password:  user.Password,
		Email:     user.Email,
		Type:      user.Type,
		CreatedAt: time.Unix(user.CreatedAt, 0),
		UpdatedAt: time.Unix(user.UpdatedAt, 0),
	}

}

func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte("tesodev")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userName"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
