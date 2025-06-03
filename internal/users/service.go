package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, request CreateUserRequest) (UserResponse, error)
	Update(ctx context.Context, id string, request UpdateUserRequest)(UserResponse, error)
	Delete(ctx context.Context, id string) (error)
	ChangePassword(ctx context.Context, request ChangePasswordRequest) (UserResponse, error)
	VerifyUser(ctx context.Context, request VerifyUserRequest) (UserResponse, error)
	FindByUsername(ctx context.Context, username string) (UserResponse, error)
	FindByID(ctx context.Context, id string) (UserResponse, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (UserResponse, error)
	FindAll(ctx context.Context) ([]UserResponse, error)
	FindByGroupId(ctx context.Context, groupID string) ([]UserResponse, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo}
}

func mapUserToResponse(user User) UserResponse {
	passwordchangedat := user.PasswordChangedAt.Format(time.RFC3339)
	response := UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Surname:  user.Surname,
		PhoneNumber: user.PhoneNumber,
		Email: user.Email,
		Username: user.Username,
		IsVerified: user.IsVerified,
		IsActive: user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		PasswordChangedAt: &passwordchangedat,				
	}
	return response
}

func(s *userService) Create(ctx context.Context, request CreateUserRequest) (UserResponse, error) {
	// Validate the request
	if request.Username == "" || request.Password == "" || 
		request.PhoneNumber == "" || request.Email == "" {
		return UserResponse{} , ErrInvalidRequest
	}

	// Check if the username is already taken
	_, err := s.repo.FindByUsername(ctx, request.Username)	
	if err == nil {		
		return UserResponse{}, ErrUsernameTaken
	}
	_, err = s.repo.FindByPhoneNumber(ctx, request.PhoneNumber)	
	if err == nil {		
		return UserResponse{}, ErrPhoneNumberTaken
	}
	if err != ErrUserNotFound {		
		return UserResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        return UserResponse{}, err
    }
	// Create a new User instance
	now := time.Now()
	user := User{
		ID:       uuid.New(),
		Name:     request.Name,
		Surname:  request.Surname,
		PhoneNumber: request.PhoneNumber,
		Email: request.Email,
		Username: request.Username,
		Password: string(hashedPassword),
		IsVerified: false,
		IsActive: true,		
		CreatedAt: now,
		PasswordChangedAt: &now,
	}

	// Insert the user into the database
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}

	//Create a validation string store in redis and send it to the user via email or sms

	return mapUserToResponse(createdUser), nil
}

func(s *userService) Update(ctx context.Context,id string,  request UpdateUserRequest)(UserResponse, error){
	if id == "" || request.Username == "" || request.PhoneNumber == "" {
		return UserResponse{} , ErrInvalidRequest
	}

	existingUser, err := s.repo.FindByID(ctx, uuid.MustParse(id))
	if err != nil && !errors.Is(err, ErrUserNotFound){
		return UserResponse{}, fmt.Errorf("something went wrong getting the user")
	}
	if errors.Is(err, ErrUserNotFound) {
		return UserResponse{}, err
	}
	if !existingUser.IsActive {
		return UserResponse{}, ErrInactiveUser
	}
	now := time.Now()
	user := User{
		ID:       uuid.MustParse(id),
		Name:     request.Name,
		Surname:  request.Surname,
		PhoneNumber: request.PhoneNumber,
		Email: request.Email,
		Username: request.Username,		
		IsActive: request.IsActive,		
		PasswordChangedAt: &now,
	}

	response, err := s.repo.Update(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}
	response.IsVerified = existingUser.IsVerified
	response.CreatedAt = existingUser.CreatedAt
	return mapUserToResponse(response), nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	existingUser, err := s.repo.FindByID(ctx, uuid.MustParse(id))
	if err != nil && !errors.Is(err, ErrUserNotFound){
		return  fmt.Errorf("something went wrong getting the user")
	}
	if errors.Is(err, ErrUserNotFound) {
		return err
	}
	if !existingUser.IsActive {
		return ErrInactiveUser
	}

	err = s.repo.Delete(ctx, parsedID)
	if err != nil {		
		return err
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, request ChangePasswordRequest) (UserResponse, error) {
	if request.ID == "" || request.Password == "" {
		return UserResponse{}, ErrInvalidRequest
	}

	existingUser, err := s.repo.FindByID(ctx, uuid.MustParse(request.ID))
	if err != nil && !errors.Is(err, ErrUserNotFound){
		return UserResponse{}, fmt.Errorf("something went wrong getting the user")
	}
	if errors.Is(err, ErrUserNotFound) {
		return UserResponse{}, err
	}
	if !existingUser.IsActive {
		return UserResponse{}, ErrInactiveUser
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}

	request.Password = string(hashedPassword)

	response, err := s.repo.ChangePassword(ctx, request)
	if err != nil {
		return UserResponse{}, err
	}
	response.IsVerified = existingUser.IsVerified
	response.CreatedAt = existingUser.CreatedAt
	return mapUserToResponse(response), nil
}

func (s *userService) VerifyUser(ctx context.Context, request VerifyUserRequest) (UserResponse, error) {
	if request.ID == "" || request.ValidationString == "" {
		return UserResponse{}, ErrInvalidRequest
	}

	existingUser, err := s.repo.FindByID(ctx, uuid.MustParse(request.ID))
	if err != nil && !errors.Is(err, ErrUserNotFound){
		return UserResponse{}, fmt.Errorf("something went wrong getting the user")
	}
	if errors.Is(err, ErrUserNotFound) {
		return UserResponse{}, err
	}
	if !existingUser.IsActive {
		return UserResponse{}, ErrInactiveUser
	}
	//Validate the validation string


	err = s.repo.VerifyUser(ctx, uuid.MustParse(request.ID))
	if err != nil {
		return UserResponse{}, err
	}

	response, err := s.repo.FindByID(ctx, uuid.MustParse(request.ID))
	if err != nil {
		return UserResponse{}, err
	}

	return mapUserToResponse(response), nil
}

func (s *userService) FindByUsername(ctx context.Context, username string) (UserResponse, error) {
	if username == "" {
		return UserResponse{}, ErrInvalidRequest
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrUserNotFound
		}
		return UserResponse{}, err
	}

	return mapUserToResponse(user), nil
}

func (s *userService) FindByPhoneNumber(ctx context.Context, phoneNumber string) (UserResponse, error) {
	if phoneNumber == "" {
		return UserResponse{}, ErrInvalidRequest
	}

	user, err := s.repo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrUserNotFound
		}
		return UserResponse{}, err
	}

	return mapUserToResponse(user), nil
}

func (s *userService) FindByID(ctx context.Context, id string) (UserResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return UserResponse{}, ErrInvalidID
	}

	user, err := s.repo.FindByID(ctx, parsedID) // Usa l'ID validat
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrUserNotFound
		}
		return UserResponse{}, err
	}

	return mapUserToResponse(user), nil
}

func (s *userService) FindAll(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, mapUserToResponse(user))
	}

	return userResponses, nil
}

func (s *userService) FindByGroupId(ctx context.Context, groupID string) ([]UserResponse, error) {
	parsedID, err := uuid.Parse(groupID)
	if err != nil {
		return nil, ErrInvalidID
	}

	users, err := s.repo.FindByGroupId(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, mapUserToResponse(user))
	}

	return userResponses, nil
}