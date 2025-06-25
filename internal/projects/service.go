package projects

import (
	"context"
	"fmt"
	"orkestra-api/internal/customers"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProjectService interface {
	Create(ctx context.Context, request ProjectRequest) (Project, error)
	Update(ctx context.Context, id string, request ProjectRequest) (Project, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (Project, error)		
	FindAllByUserID(ctx context.Context, id string) ([]Project, error)
	FindBetweenDates(ctx context.Context, startDate, endDate string) ([]Project, error)	
	FindCalendarBetweenDatesByUserID(ctx context.Context,userID, startDate, endDate string)([]ProjectCalendarResponse, error)
}

type projectService struct {
	repo ProjectRepository
	customerService customers.CustomerService
}

func NewProjectService(repo ProjectRepository, customerService customers.CustomerService)ProjectService{
	return &projectService{repo, customerService}
}

func(s *projectService) Create(ctx context.Context, request ProjectRequest) (Project, error){
	project, err := createModelFromRequest(request)
	if err != nil {
		return Project{}, err
	}
	ret, err := s.repo.Create(ctx, project)
	if err != nil {
		return Project{}, err
	}
	return ret, nil
}
func(s *projectService) Update(ctx context.Context, id string, request ProjectRequest) (Project, error){
	project, err := createModelFromRequest(request)
	if err != nil {
		return Project{}, err
	}
	projectUUID, err := uuid.Parse(id)
	if err != nil {
		return Project{}, ErrInvalidID
	}
	project.ID = projectUUID
	ret, err := s.repo.Update(ctx, project)
	if err != nil {
		return Project{}, err
	}
	return ret, nil
}
func(s *projectService) Delete(ctx context.Context, id string) error{
	projectUUID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, projectUUID)
}
func(s *projectService) FindById(ctx context.Context, id string) (Project, error){
	projectUUID, err := uuid.Parse(id)
	if err != nil {
		return Project{}, ErrInvalidID
	}
	return s.repo.FindById(ctx, projectUUID)
}

func(s *projectService) FindAllByUserID(ctx context.Context, userID string) ([]Project, error){
	customer, err := s.customerService.FindCustomerByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding customer by user ID: %w", err)
	}
	if customer.ID == uuid.Nil {
		return s.findAll(ctx)
	}
	return s.findAllByUserID(ctx, userID)
}

func(s *projectService) findAll(ctx context.Context) ([]Project, error){
	return s.repo.FindAll(ctx)
}

func(s *projectService) findAllByUserID(ctx context.Context, userID string) ([]Project, error){
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindAllByUserID(ctx, userUUID)
}

func(s *projectService) FindBetweenDates(ctx context.Context, startDate, endDate string) ([]Project, error){
	layout := "2006-01-02" 

	pStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	pEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	return s.repo.FindBetweenDates(ctx, &pStartDate, &pEndDate)
}

func(s *projectService)FindCalendarBetweenDatesByUserID(ctx context.Context, userID, startDate, endDate string)([]ProjectCalendarResponse, error){
	customer, err := s.customerService.FindCustomerByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding customer by user ID: %w", err)
	}
	if customer.ID == uuid.Nil {
		return s.findCalendarBetweenDates(ctx, startDate, endDate)
	}
	return s.findCalendarBetweenDatesByUserID(ctx, userID, startDate, endDate)
}

func(s *projectService)findCalendarBetweenDates(ctx context.Context, startDate, endDate string)([]ProjectCalendarResponse, error){
	layout := "2006-01-02" 

	pStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	pEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	return s.repo.FindCalendarBetweenDates(ctx, &pStartDate, &pEndDate)
}

func(s *projectService)findCalendarBetweenDatesByUserID(ctx context.Context, userID, startDate, endDate string)([]ProjectCalendarResponse, error){
	layout := "2006-01-02" 

	pStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	pEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindCalendarBetweenDatesByUserID(ctx, userUUID, &pStartDate, &pEndDate)
}

func createModelFromRequest(request ProjectRequest)(Project, error){
if 	request.Color == "" ||
		request.CustomerID == "" ||
		request.Description == "" ||
		request.StartDate == "" ||
		request.EndDate == "" {
			return Project{}, ErrInvalidRequest
		}
	layout := time.RFC3339 

	startDate, err := time.Parse(layout, request.StartDate)
	if err != nil {
		return Project{}, ErrInvalidDate
	}
	endDate, err := time.Parse(layout, request.EndDate)
	if err != nil {
		return Project{}, ErrInvalidDate
	}

	customerUUID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Project{}, ErrInvalidID
	}

	amount, err := decimal.NewFromString(request.Amount)
	if err != nil {
		return Project{}, ErrInvalidRequest
	}
	estimatedCost, err := decimal.NewFromString(request.EstimatedCost)
	if err != nil {
		return Project{}, ErrInvalidRequest
	}

	project := Project{
		ID: uuid.New(),
		Description: request.Description,
		Color: request.Color,
		StartDate: &startDate,
		EndDate: &endDate,
		CustomerID: customerUUID,
		Amount: amount,
		EstimatedCost: estimatedCost,
	}
	return project, nil
}