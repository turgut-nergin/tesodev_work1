package response

import "time"

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
