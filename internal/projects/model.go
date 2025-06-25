package projects

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Project struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	StartDate   *time.Time `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date" db:"end_date"`
	Color       string `json:"color" db:"color"`
	CustomerID  uuid.UUID `json:"customer_id" db:"customer_id"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	EstimatedCost decimal.Decimal `json:"estimated_cost" db:"estimated_cost"`
}