package tasks

type TaskRequest struct {
	Description string       `json:"description" binding:"required"`
	Notes       string       `json:"notes"`
	UserID      string       `json:"user_id" binding:"required"`
	Status      TaskStatus   `json:"status" binding:"required,oneof=Pending ToDo InProgress Done"`
	Priority    TaskPriority `json:"priority" binding:"required,oneof=A B C D"`
	ProjectID   string       `json:"project_id" binding:"required"`
	StartDate   *string      `json:"start_date,omitempty"`
	EndDate     *string      `json:"end_date,omitempty"`
}
