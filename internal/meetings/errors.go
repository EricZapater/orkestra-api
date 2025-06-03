package meetings

import "errors"

var (
	ErrInvalidID            = errors.New("invalid meeting ID")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrMeetingNotFound      = errors.New("meeting not found")
	ErrMeetingAlreadyExists = errors.New("meeting already exists")
)