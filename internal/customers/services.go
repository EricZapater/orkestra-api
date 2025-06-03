package customers

import (
	"context"

	"github.com/google/uuid"
)

type CustomerService interface {
	Create(ctx context.Context, request CustomerRequest)(Customer, error)
	Update(ctx context.Context,id string, request CustomerRequest)(Customer, error)
	Delete(ctx context.Context, id string)(error)
	FindByID(ctx context.Context, id string)(Customer, error)
	FindAll(ctx context.Context)([]Customer, error)
	AddUserToCustomer(ctx context.Context, request UserCustomerRequest) (error)
	RemoveUserFromCustomer(ctx context.Context, request UserCustomerRequest) (error)
}

type customerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService{
	return &customerService{repo}
}

func(s *customerService) Create(ctx context.Context, request CustomerRequest)(Customer, error){
	if request.ComercialName == "" || request.PhoneNumber == "" || request.VatNumber == "" {
		return Customer{}, ErrInvalidRequest
	}
	customer := Customer{
		ID: uuid.New(),
		ComercialName: request.ComercialName,
		PhoneNumber: request.PhoneNumber,
		VatNumber: request.VatNumber,
	}

	return s.repo.Create(ctx, customer)
}

func(s *customerService) Update(ctx context.Context, id string, request CustomerRequest)(Customer, error){
	customerID, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, ErrInvalidID
	}
	if request.ComercialName == "" || request.PhoneNumber == "" || request.VatNumber == "" {
		return Customer{}, ErrInvalidRequest
	}
	customer := Customer{
		ID: customerID,
		ComercialName: request.ComercialName,
		PhoneNumber: request.PhoneNumber,
		VatNumber: request.VatNumber,
	}

	return s.repo.Update(ctx, customer)	
}

func(s *customerService) Delete(ctx context.Context, id string)(error){
	customerID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, customerID)
}
func(s *customerService) FindByID(ctx context.Context, id string)(Customer, error){
	customerID, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, ErrInvalidID
	}
	return s.repo.FindById(ctx, customerID)
}
func(s *customerService) FindAll(ctx context.Context)([]Customer, error){
	return s.repo.FindAll(ctx)
}

func(s *customerService) AddUserToCustomer(ctx context.Context, request UserCustomerRequest) (error){
	id := uuid.New()
	customerUUID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return  ErrInvalidID
	}
	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		return  ErrInvalidID
	}
	return s.repo.AddUserToCustomer(ctx, id, customerUUID, userUUID)
}
func(s *customerService) RemoveUserFromCustomer(ctx context.Context, request UserCustomerRequest) (error){
	customerUUID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return  ErrInvalidID
	}
	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		return  ErrInvalidID
	}
	return s.repo.RemoveUserFromCustomer(ctx, customerUUID, userUUID)
}