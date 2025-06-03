package auth

import (
	"context"
	"orkestra-api/internal/users"

	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"golang.org/x/crypto/bcrypt"
)



type AuthService interface {
    Login(ctx context.Context, req LoginRequest) (string, users.User, time.Time, error)
    ValidateUser(username, password string) (users.User, error)
}

type authService struct {
    userRepo users.UserRepository
    jwtMiddleware *jwt.GinJWTMiddleware
}

func NewAuthService(userRepo users.UserRepository,  jwtMiddleware *jwt.GinJWTMiddleware) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtMiddleware: jwtMiddleware,
    }
}

// Login verifica les credencials i retorna un token JWT si són vàlides
func (s *authService) Login(ctx context.Context, req LoginRequest) (string, users.User, time.Time, error) {
    // Validar les credencials
    user, err := s.ValidateUser(req.Username, req.Password)
    if err != nil {
        return "", users.User{}, time.Time{}, err
    }
    id := user.ID.String()
    // Generar token JWT
    token, expire, err := s.jwtMiddleware.TokenGenerator(id)
    if err != nil {
        return "", users.User{}, time.Time{}, err
    }
   
    user.Password = "" // No retornar la contrasenya en la resposta
    return token, user,  expire, nil
}

// ValidateUser verifica si les credencials són vàlides i retorna l'ID de l'usuari
func (s *authService) ValidateUser(username, password string) (users.User, error) {
    // Obtenir l'usuari per nom d'usuari
    user, err := s.userRepo.FindByUsername(context.Background(), username)

    if err != nil {
        return users.User{}, ErrUserNotFound
    }
    
    // Verificar que l'usuari estigui actiu
    if !user.IsActive {
        return users.User{}, users.ErrInactiveUser
    }
    
    // Verificar la contrasenya
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return users.User{}, ErrInvalidCredentials
    }
    
    // Retornar l'ID de l'usuari com a identificador principal
    return user, nil
}