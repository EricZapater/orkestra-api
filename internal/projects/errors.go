package projects

import "errors"

var (
	ErrInvalidID      = errors.New("invalid project ID")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidDate    = errors.New("invalid date")
	ErrProjectNotFound   = errors.New("project not found")
	ErrOperatorDatesOutOfProjectRange = errors.New("operator's start and end dates must be within the project's date range")
)