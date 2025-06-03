package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {	
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uuid.UUID) (error)
	ChangePassword(ctx context.Context, request ChangePasswordRequest) (User, error)
	VerifyUser(ctx context.Context, id uuid.UUID) (error)
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)	
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByGroupId(ctx context.Context, groupId uuid.UUID) ([]User,error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user User) (User, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, surname, phone_number, email, username, password, is_verified, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Name, user.Surname, user.PhoneNumber, user.Email, user.Username, user.Password, user.IsVerified, user.IsActive,
    )
    if err != nil {
        return User{}, fmt.Errorf("error inserting user: %w", err)
    }
    return user, nil
}

func(r *userRepository) Update(ctx context.Context, user User) (User, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET name = $1, surname = $2, phone_number = $3, email = $4, username = $5, is_active = $6
		WHERE id = $7`,
		user.Name, user.Surname, user.PhoneNumber, user.Email, user.Username, user.IsActive, user.ID,
	)
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %w", err)
	}
	return user, nil
}

func(r *userRepository) Delete(ctx context.Context, id uuid.UUID) (error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET is_active = false
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func(r *userRepository) ChangePassword(ctx context.Context, request ChangePasswordRequest) (User, error){
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET password = $1, password_changed_at = now()
		WHERE id = $2`,
		request.Password, request.ID)
	if err != nil {
		return User{}, fmt.Errorf("error changing password: %w", err)
	}
	strId, err  := uuid.Parse(request.ID)
	if err != nil {
		return User{}, fmt.Errorf("error parsing id: %w", err)
	}
	user, err := r.FindByID(ctx, strId)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func(r *userRepository) VerifyUser(ctx context.Context, id uuid.UUID) (error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET is_verified = true
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("error verifying user: %w", err)
	}
	return nil
}

func(r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (User, error){
	var user User
	row := r.db.QueryRowContext(ctx, `SELECT id, name, surname, phone_number, email, username, password, is_verified, is_active, created_at, password_changed_at FROM users WHERE id = $1`, id)
	
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}else if err != nil {
		return User{}, fmt.Errorf("error getting user: %w", err)
	}
	if !user.IsActive {
		return User{}, ErrInactiveUser
	}
	return user, nil
}

func(r *userRepository) FindByUsername(ctx context.Context, username string) (User, error)	{
	var user User
	row := r.db.QueryRowContext(ctx, `SELECT id, name, surname, phone_number, email, username, password, is_verified, is_active, created_at, password_changed_at FROM users WHERE username = $1`, username)
	
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}else if err != nil {
		return User{}, fmt.Errorf("error getting user: %w", err)
	}
	if !user.IsActive {
		return User{}, ErrInactiveUser
	}
	return user, nil
}

func(r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)	{
	var user User
	row := r.db.QueryRowContext(ctx, `SELECT id, name, surname, phone_number, email, username, password, is_verified, is_active, created_at, password_changed_at FROM users WHERE phone_number = $1`, phoneNumber)
	
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}else if err != nil {
		return User{}, fmt.Errorf("error getting user: %w", err)
	}
	if !user.IsActive {
		return User{}, ErrInactiveUser
	}
	return user, nil
}

func(r *userRepository) FindAll(ctx context.Context) ([]User, error){
	var users []User
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, surname, phone_number, email, username, password, is_verified, is_active, created_at, password_changed_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func(r *userRepository) FindByGroupId(ctx context.Context, groupId uuid.UUID) ([]User,error) {
	var users []User
	rows, err := r.db.QueryContext(ctx, 
		`SELECT u.id, u.name, u.surname, u.phone_number, u.email, u.username, u.password, u.is_verified, u.is_active, u.created_at, u.password_changed_at 
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = $1`, groupId)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}