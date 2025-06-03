package projects

import "errors"

var (
	ErrInvalidID      = errors.New("invalid project ID")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidDate    = errors.New("invalid date")
)