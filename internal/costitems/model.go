package costitems

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CostItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Amount    decimal.Decimal `json:"amount" db:"decimal"`
	ShortDescription string `json:"short_description" db:"short_description"`
	Notes string `json:"notes" db:"notes"`
	Date *time.Time `json:"date" db:"date"`
}