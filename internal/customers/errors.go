package customers

import "errors"

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidID        = errors.New("invalid customer ID")
	ErrCustomerVatNumberTaken   = errors.New("customer vat number already taken")
	ErrInvalidRequest   = errors.New("invalid request")
)