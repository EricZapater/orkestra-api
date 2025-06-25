package tasks

import "errors"

var (
	ErrInvalidID      = errors.New("invalid task ID")
	ErrInvalidRequest = errors.New("invalid task request")
	ErrInvalidStatus = errors.New("invalid task status")
	ErrInvalidPriority = errors.New("invalid task priority")
	ErrInvalidDate    = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrTaskNotFound   = errors.New("task not found")
)