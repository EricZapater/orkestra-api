package menus

import "github.com/google/uuid"

type Menu struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Label       string    `json:"label" db:"label"`
	Icon        string    `json:"icon" db:"icon"`
	Route       string    `json:"route" db:"route"`
	ParentID    uuid.UUID `json:"parent_id" db:"parent_id"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	IsSeparator bool      `json:"is_separator" db:"is_separator"`
}