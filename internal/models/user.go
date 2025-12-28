package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Name      string    `json:"name,omitempty" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the request to create/update a user
type CreateUserRequest struct {
	Phone string `json:"phone" binding:"required,min=10,max=15"`
	Name  string `json:"name" binding:"max=255"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required,max=255"`
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.Phone == "" {
		return ErrInvalidPhone
	}
	if len(u.Phone) < 10 || len(u.Phone) > 15 {
		return ErrInvalidPhone
	}
	return nil
}
