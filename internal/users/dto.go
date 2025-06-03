package users

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Username    string `json:"username" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type ChangePasswordRequest struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUserRequest struct {
	ID               string `json:"id" binding:"required"`
	ValidationString string `json:"validation_string" binding:"required"`
}

type UserResponse struct {
	ID                string  `json:"id" db:"id"`
	Name              string  `json:"name" db:"name"`
	Surname           string  `json:"surname" db:"surname"`
	PhoneNumber       string  `json:"phone_number" db:"phone_number"`
	Email             string  `json:"email" db:"email"`
	Username          string  `json:"username" db:"username"`
	Password          string  `json:"password" db:"password"`
	IsVerified        bool    `json:"is_verified" db:"is_verified"`
	IsActive          bool    `json:"is_active" db:"is_active"`
	CreatedAt         string  `json:"created_at" db:"created_at"`
	PasswordChangedAt *string `json:"password_changed_at" db:"password_changed_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}