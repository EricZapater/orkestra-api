package costitems

import "errors"

var (
	ErrInvalidID      = errors.New("invalid task ID")
	ErrInvalidRequest = errors.New("invalid task request")	
	ErrInvalidDate    = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrProjectNotFound   = errors.New("project not found")
)