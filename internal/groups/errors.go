package groups

import "errors"

var (
	ErrGroupNotFound  = errors.New("group not found")
	ErrInvalidID      = errors.New("invalid group ID")
	ErrGroupNameTaken = errors.New("group name already taken")
	ErrInvalidRequest = errors.New("invalid request")
)