package searches

type SearchRequest struct {
	GroupIDS []string `json:"group_ids" binding:"required,min=1,dive,uuid"`
	Text     string   `json:"text" binding:"required"`
}