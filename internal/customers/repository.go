package customers

import (
	"context"
	"database/sql"
	"fmt"
	"orkestra-api/internal/users"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) (Customer, error)
	Update(ctx context.Context, customer Customer) (Customer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (Customer, error)
	FindAll(ctx context.Context)([]Customer, error)
	AddUserToCustomer(ctx context.Context, id, customerID, userID uuid.UUID) (error)
	RemoveUserFromCustomer(ctx context.Context, customerID, userID uuid.UUID) (error)
	FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (Customer, error)
	FindUsersByCustomerID(ctx context.Context, customerID uuid.UUID) ([]users.User, error)
}

type customerRepository struct{
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{
		db:db,
	}
}

func(r *customerRepository) Create(ctx context.Context, customer Customer) (Customer, error){
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO customers(id, comercial_name, vat_number, phone_number)
	VALUES($1, $2, $3, $4)
	`,
	customer.ID, customer.ComercialName, customer.VatNumber, customer.PhoneNumber,
	)
	if err != nil {
		return Customer{}, err
	}
	return customer, err
}

func(r *customerRepository) Update(ctx context.Context, customer Customer) (Customer, error){
	_, err := r.db.ExecContext(ctx, `
		UPDATE customers
		set comercial_name = $1,
		vat_number = $2,
		phone_number = $3
		WHERE id = $4`,
		customer.ComercialName, customer.VatNumber, customer.PhoneNumber, customer.ID,
	)
	if err != nil {
		return Customer{}, err
	}
	return customer, err
}
func(r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error{
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1`, id)
	return err
}
func(r *customerRepository) FindById(ctx context.Context, id uuid.UUID) (Customer, error){
	var customer Customer
	err := r.db.QueryRowContext(ctx, `SELECT id, comercial_name, vat_number, phone_number FROM customers WHERE id = $1`, id,
).Scan(&customer.ID, &customer.ComercialName, &customer.VatNumber, &customer.PhoneNumber)
if err != nil {
	return Customer{}, err
}
return customer, nil
}

func(r *customerRepository) FindAll(ctx context.Context)([]Customer, error){
	var customers []Customer
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, comercial_name, vat_number, phone_number FROM customers 
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.ComercialName, &customer.VatNumber, &customer.PhoneNumber); err != nil{
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
func(r *customerRepository) AddUserToCustomer(ctx context.Context, id, customerID, userID uuid.UUID) (error){
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO customer_users(id, customer_id, user_id)VALUES($1, $2, $3)
	`, id, customerID, userID)
	if err != nil {
		return err
	}
	return nil
}
func(r *customerRepository) RemoveUserFromCustomer(ctx context.Context, customerID, userID uuid.UUID) (error){
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM customer_users WHERE customer_id = $1 AND user_id = $2
	`, customerID, userID)
	if err != nil {
		return err
	}
	return nil
}
func(r *customerRepository) FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (Customer, error){
	var customer Customer
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.comercial_name, c.vat_number, c.phone_number 
		FROM customers c
		INNER JOIN customer_users cu ON c.id = cu.customer_id
		WHERE cu.user_id = $1`, userID,
).Scan(&customer.ID, &customer.ComercialName, &customer.VatNumber, &customer.PhoneNumber)
if err != nil {
	return Customer{}, err
}
return customer, nil
}
func(r *customerRepository) FindUsersByCustomerID(ctx context.Context, customerID uuid.UUID) ([]users.User, error){

	var usersList []users.User
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.name, u.surname, u.phone_number, u.email, u.username, u.password, u.is_verified, u.is_active, u.created_at, u.password_changed_at, u.profile_id 
		FROM users u
			INNER JOIN customer_users cu ON u.id = cu.user_id
		WHERE cu.customer_id = $1
		`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user users.User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.PasswordChangedAt, &user.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}
