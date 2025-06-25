package costitems

import (
	"context"
	"orkestra-api/internal/projects"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CostItemService interface {
	Create(ctx context.Context, request *CostItemRequest) (CostItem, error)
	Update(ctx context.Context, id string, request *CostItemRequest) (CostItem, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (CostItem, error)
	FindByProjectID(ctx context.Context, id string) ([]CostItem, error)
	FindAll(ctx context.Context) ([]CostItem, error)
}

type costItemService struct {
	repo CostItemRepository
	projects projects.ProjectService
}

func NewCostItemService(repo CostItemRepository, projects projects.ProjectService)CostItemService{
	return &costItemService{
		repo: repo,
		projects: projects,
	}
}

func(s *costItemService) Create(ctx context.Context, request *CostItemRequest) (CostItem, error){
	_, err := s.projects.FindById(ctx, request.ProjectID)
	if err != nil {
		return CostItem{}, ErrProjectNotFound
	}

	costItem,err := createModelFromRequest(request)
	if err != nil{
		return CostItem{}, err
	}
	
	return s.repo.Create(ctx, costItem)
}
func(s *costItemService) Update(ctx context.Context, id string, request *CostItemRequest) (CostItem, error){
	_, err := s.projects.FindById(ctx, request.ProjectID)
	if err != nil {
		return CostItem{}, ErrProjectNotFound
	}

	costItem,err := createModelFromRequest(request)
	if err != nil{
		return CostItem{}, err
	}

	costItemID, err := uuid.Parse(id)
	if err != nil{
		return CostItem{}, ErrInvalidID
	}
	costItem.ID = costItemID
	return s.repo.Update(ctx, costItem)
}

func(s *costItemService) Delete(ctx context.Context, id string) error{
	costItemID, err := uuid.Parse(id)
	if err != nil{
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, costItemID)
}

func(s *costItemService) FindByID(ctx context.Context, id string) (CostItem, error){
	costItemID, err := uuid.Parse(id)
	if err != nil{
		return CostItem{}, ErrInvalidID
	}
	return s.repo.FindByID(ctx, costItemID)
}
func(s *costItemService) FindByProjectID(ctx context.Context, id string) ([]CostItem, error){
		_, err := s.projects.FindById(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindByProjectID(ctx, projectID)
}
func(s *costItemService) FindAll(ctx context.Context) ([]CostItem, error){
	return s.repo.FindAll(ctx)
}

func createModelFromRequest(request *CostItemRequest)(CostItem, error){
	if request.Amount == "" || request.ProjectID == "" || request.ShortDescription == "" {
		return CostItem{}, ErrInvalidRequest
	}
	costItemID := uuid.New()

	projectID, err := uuid.Parse(request.ProjectID)
	if err != nil {
		return CostItem{}, ErrInvalidID
	}
	layout := time.RFC3339
	var date time.Time
	if len(request.Date) > 0 {
		date, err = time.Parse(layout, request.Date)
		if err != nil {
			return CostItem{}, ErrInvalidDate
		}
	}
	amount, err := decimal.NewFromString(request.Amount)
	if err != nil {
		return CostItem{}, ErrInvalidRequest
	}
	costItem := CostItem{
		ID: costItemID,
		ProjectID: projectID,
		Amount: amount,
		ShortDescription: request.ShortDescription,
		Notes: request.Notes,
		Date: &date,
	}
	return costItem, nil
}