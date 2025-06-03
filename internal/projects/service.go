package projects

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, request ProjectRequest) (Project, error)
	Update(ctx context.Context, id string, request ProjectRequest) (Project, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (Project, error)
	FindAll(ctx context.Context) ([]Project, error)
	FindBetweenDates(ctx context.Context, startDate, endDate string) ([]Project, error)
}

type projectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository)ProjectService{
	return &projectService{repo}
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

func(s *projectService) FindAll(ctx context.Context) ([]Project, error){
	return s.repo.FindAll(ctx)
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

func createModelFromRequest(request ProjectRequest)(Project, error){
if 	request.Color == "" ||
		request.CustomerID == "" ||
		request.Description == "" ||
		request.StartDate == "" ||
		request.EndDate == "" {
			return Project{}, ErrInvalidRequest
		}
	layout := "2006-01-02" 

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

	project := Project{
		ID: uuid.New(),
		Description: request.Description,
		Color: request.Color,
		StartDate: &startDate,
		EndDate: &endDate,
		CustomerID: customerUUID,
	}
	return project, nil
}