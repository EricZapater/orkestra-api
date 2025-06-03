package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Surname  string    `json:"surname" db:"surname"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email		   string `json:"email" db:"email"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	IsVerified bool `json:"is_verified" db:"is_verified"`
	IsActive bool `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	PasswordChangedAt *time.Time `json:"password_changed_at" db:"password_changed_at"`
}