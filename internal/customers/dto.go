package customers

type CustomerRequest struct {
	ComercialName string `json:"comercial_name" binding:"required"`
	VatNumber     string `json:"vat_number" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
}

type UserCustomerRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}