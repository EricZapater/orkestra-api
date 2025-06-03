package groups

import "github.com/google/uuid"

type Group struct {
	ID uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
}