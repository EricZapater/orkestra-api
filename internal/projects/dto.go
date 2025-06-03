package projects

type ProjectRequest struct {
	Description string `json:"description" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	Color       string `json:"color" binding:"required"`
	CustomerID  string `json:"customer_id" binding:"required"`
}