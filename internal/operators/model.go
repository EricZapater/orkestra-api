package operators

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Operator struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Surname string    `json:"surname" db:"surname"`
	Cost    decimal.Decimal `json:"cost" db:"cost"`
	Color   string    `json:"color" db:"color"`
}