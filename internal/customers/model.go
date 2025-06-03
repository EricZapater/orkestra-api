package customers

import (
	"orkestra-api/internal/users"

	"github.com/google/uuid"
)

type Customer struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ComercialName string    `json:"comercial_name" db:"comercial_name"`
	VatNumber     string    `json:"vat_number" db:"vat_number"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Users *[]users.User `json:"users"`
}