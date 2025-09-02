package profiles

type ProfileRequest struct {
	Name string `json:"name" binding:"required"`
}