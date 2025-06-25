package tasks

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID       `json:"id" db:"required"`
	Description string       `json:"description" db:"description"`
	Notes       string       `json:"notes" db:"notes"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	ProjectID   uuid.UUID    `json:"project_id" db:"project_id"`
	StartDate *time.Time `json:"start_date,omitempty" db:"start_date"`
 	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"`
}