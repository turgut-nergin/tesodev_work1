package lib

import (
	"time"

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
