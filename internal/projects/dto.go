package projects

type ProjectRequest struct {
	Description   string `json:"description" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	Color         string `json:"color" binding:"required"`
	CustomerID    string `json:"customer_id" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	EstimatedCost string `json:"estimated_cost" binding:"required"`
}

type ProjectCalendarResponse struct {
	ID        string `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Color     string `json:"color" binding:"required"`
}

type OperatorToProjectRequest struct {
	OperatorID        string `json:"operator_id" binding:"required"`
	ProjectID         string `json:"project_id" binding:"required"`
	Cost              string `json:"cost" `
	DedicationPercent string `json:"dedication_percentage"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date" `
}

type CostItemRequest struct {
	ProjectID        string `json:"project_id" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Notes            string `json:"notes" binding:"required"`
	Date             string `json:"date" binding:"required"`
}