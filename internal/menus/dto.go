package menus

type CreateMenuRequest struct {
	Label       string    `json:"label" binding:"required"`
	Icon        string    `json:"icon" binding:"required"`
	Route       string    `json:"route" binding:"required"`
	ParentID    string `json:"parent_id" binding:"required"`
	SortOrder   int       `json:"sort_order" binding:"required"`
	IsSeparator bool      `json:"is_separator" binding:"required"`
}