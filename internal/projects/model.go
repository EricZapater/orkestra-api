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

type CostItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Amount    decimal.Decimal `json:"amount" db:"decimal"`
	ShortDescription string `json:"short_description" db:"short_description"`
	Notes string `json:"notes" db:"notes"`
	Date *time.Time `json:"date" db:"date"`
}

type OperatorToProject struct {
	ID				uuid.UUID `json:"id" db:"id"`
	OperatorID uuid.UUID `json:"operator_id" binding:"required"`
	ProjectID  uuid.UUID `json:"project_id" binding:"required"`
	Cost decimal.Decimal `json:"cost" binding:"required"`
	DedicationPercent decimal.Decimal `json:"dedication_percent" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}