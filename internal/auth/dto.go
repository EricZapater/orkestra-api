package auth

import (
	"orkestra-api/internal/users"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
	User   users.User   `json:"user"`	
}