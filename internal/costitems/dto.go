package costitems

type CostItemRequest struct {
	ProjectID        string `json:"project_id" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Notes            string `json:"notes" binding:"required"`
	Date             string `json:"date" binding:"required"`
}