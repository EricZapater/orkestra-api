package auth

import (
	"net/http"
	"orkestra-api/internal/users"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService    AuthService
    jwtMiddleware *jwt.GinJWTMiddleware
}

func NewAuthHandler(authService AuthService, jwtMiddleware *jwt.GinJWTMiddleware) *AuthHandler {
    return &AuthHandler{
        authService:    authService,
        jwtMiddleware: jwtMiddleware,
    }
}

// Login processa una petició de login i retorna un token JWT
func (h *AuthHandler) Login(c *gin.Context) {
    var loginRequest LoginRequest
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request"})
        return
    }
    
    token,user, expire, err := h.authService.Login(c.Request.Context(), loginRequest)
    if err != nil {
        var statusCode int
        switch err {
        case ErrInvalidCredentials:
            statusCode = http.StatusUnauthorized
        case users.ErrInactiveUser:
            statusCode = http.StatusForbidden
        default:
            statusCode = http.StatusInternalServerError
        }
        
        c.JSON(statusCode, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, LoginResponse{
        Token:  token,
        Expire: expire.Format(time.RFC3339),
        User:   user,
    })
}

// RefreshToken utilitza el middleware JWT per refrescar el token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    h.jwtMiddleware.RefreshHandler(c)
}
