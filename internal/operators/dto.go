package operators

type OperatorRequest struct {
	Name    string `json:"name" binding:"required"`
	Surname string `json:"surname"`
	Cost    string `json:"cost" binding:"required"`
	Color   string `json:"color" binding:"required"`
}