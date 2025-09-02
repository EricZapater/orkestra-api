package operators

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperatorService interface {
	Create(ctx context.Context, request OperatorRequest)(Operator, error)
	Update(ctx context.Context, id string, request OperatorRequest)(Operator, error)
	Delete(ctx context.Context, id string)error
	FindByID(ctx context.Context, id string)(Operator, error)
	FindAll(ctx context.Context)([]Operator, error)
}

type operatorService struct {
	repo OperatorRepository
}

func NewOperatorService(repo OperatorRepository)OperatorService{
	return &operatorService{
		repo:repo,
	}
}

func(s* operatorService) Create(ctx context.Context, request OperatorRequest)(Operator, error){
	cost, err := decimal.NewFromString(request.Cost)
	if err != nil {
		return Operator{}, err
	}
	operator := Operator{
		ID: uuid.New(),
		Name: request.Name,
		Surname: request.Surname,
		Cost: cost,
		Color: request.Color,
	}
	return s.repo.Create(ctx, operator)
}

func(s* operatorService) Update(ctx context.Context, id string, request OperatorRequest)(Operator, error){
	operatorUUID, err := uuid.Parse(id)
	if err != nil {
		return Operator{}, err
	}
	cost, err := decimal.NewFromString(request.Cost)
	if err != nil {
		return Operator{}, err
	}
	operator := Operator{
		ID: operatorUUID,
		Name: request.Name,
		Surname: request.Surname,
		Cost: cost,
		Color: request.Color,
	}
	return s.repo.Update(ctx, operator)
}
func(s* operatorService) Delete(ctx context.Context, id string)error{
	operatorUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, operatorUUID)
}
func(s* operatorService) FindByID(ctx context.Context, id string)(Operator, error){
	operatorUUID, err := uuid.Parse(id)
	if err != nil {
		return Operator{}, err
	}
	return s.repo.FindByID(ctx, operatorUUID)
}
func(s* operatorService) FindAll(ctx context.Context)([]Operator, error){
	return s.repo.FindAll(ctx)
}