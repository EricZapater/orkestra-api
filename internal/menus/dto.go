package menus

type CreateMenuRequest struct {
	Label       string  `json:"label" binding:"required"`
	Icon        string  `json:"icon" binding:"required"`
	Route       string  `json:"route" binding:"required"`
	ParentID    *string `json:"parent_id"`
	SortOrder   int     `json:"sort_order" binding:"required"`
	IsSeparator bool    `json:"is_separator" default:"false"`
}

type MenuToProfileRequest struct {
	ProfileID string `json:"profile_id" binding:"required"`
	MenuID    string `json:"menu_id" binding:"required"`
}