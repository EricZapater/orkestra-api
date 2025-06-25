package menus

import "errors"

var (
	ErrInvalidParentID = errors.New("invalid parent ID")
	ErrInvalidMenuID   = errors.New("invalid menu ID")
	ErrInvalidProfileID = errors.New("invalid profile ID")
)