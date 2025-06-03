package groups

type GroupRequest struct {
	Name string `json:"name" binding:"required"`
}
