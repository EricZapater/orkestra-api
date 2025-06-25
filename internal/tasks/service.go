package tasks

import (
	"context"
	"orkestra-api/internal/projects"
	"orkestra-api/internal/users"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	Create(ctx context.Context, request *TaskRequest) (Task, error)
	Update(ctx context.Context, id string, request *TaskRequest) (Task, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Task, error)
	FindByStatus(ctx context.Context, status string) ([]Task, error)
	FindByUserID(ctx context.Context, id string) ([]Task, error)
	FindByProjectID(ctx context.Context, id string) ([]Task, error)
	FindByPriority(ctx context.Context, priority string) ([]Task, error)
	FindAll(ctx context.Context) ([]Task, error)
}

type taskService struct {
	repo TaskRepository
	users users.UserService
	projects projects.ProjectService
}

func NewTaskService(repo TaskRepository, users users.UserService, projects projects.ProjectService) TaskService {
	return &taskService{
		repo:repo,
		users:users,
		projects: projects,
	}
}

func(s *taskService) Create(ctx context.Context, request *TaskRequest) (Task, error){
	task, err := createModelFromRequest(request)
	if err != nil {
		return Task{}, err
	}
	return s.repo.Create(ctx, task)
}
func(s *taskService) Update(ctx context.Context, id string, request *TaskRequest) (Task, error){
	task, err := createModelFromRequest(request)
	if err != nil {
		return Task{}, err
	}
	taskID, err := uuid.Parse(id)
	if err != nil {
		return Task{}, ErrInvalidID
	}
	task.ID = taskID
	now := time.Now()
	if task.Status == StatusInProgress {
		task.StartDate = &now
	}
	if task.Status == StatusDone {
		task.EndDate = &now
	}
	return s.repo.Update(ctx, task)
}
func(s *taskService) Delete(ctx context.Context, id string) error{
	taskID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, taskID)
}
func(s *taskService) FindByID(ctx context.Context, id string) (Task, error){
	taskID, err := uuid.Parse(id)
	if err != nil {
		return Task{}, ErrInvalidID
	}
	return s.repo.FindById(ctx, taskID)
}
func(s *taskService) FindByStatus(ctx context.Context, status string) ([]Task, error){
	if !IsValidStatus(TaskStatus(status)){
		return nil, ErrInvalidStatus
	}
	taskStatus := TaskStatus(status)	
	return s.repo.FindByStatus(ctx, taskStatus)
}
func(s *taskService) FindByUserID(ctx context.Context, id string) ([]Task, error){
	_, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	userID, err  := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindByUserID(ctx, userID)
}
func(s *taskService) FindByProjectID(ctx context.Context, id string) ([]Task, error){
	_, err := s.projects.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	projectID, err := uuid.Parse(id)
	if err != nil{
		return nil, ErrInvalidID
	}
	return s.repo.FindByProjectID(ctx, projectID)	
}
func(s *taskService) FindByPriority(ctx context.Context, priority string) ([]Task, error){
	if !IsValidPriority(TaskPriority(priority)){
		return nil, ErrInvalidPriority
	}
	taskPriority := TaskPriority(priority)
	return s.repo.FindByPriority(ctx, taskPriority)
}
func(s *taskService) FindAll(ctx context.Context) ([]Task, error){
	return s.repo.FindAll(ctx)
}

func createModelFromRequest(request *TaskRequest)(Task, error){
if request.Description == "" || request.UserID == "" || request.Status == "" || request.Priority == "" || request.ProjectID == "" {
		return Task{}, ErrInvalidRequest
	}
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return Task{}, ErrInvalidID
	}
	projectID, err := uuid.Parse(request.ProjectID)
	if err != nil {
		return Task{}, ErrInvalidID
	}
	if !IsValidStatus(request.Status){
		return Task{}, ErrInvalidStatus
	}
	if !IsValidPriority(request.Priority){
		return Task{}, ErrInvalidPriority
	}
	layout := time.RFC3339 
	var startDate time.Time
	var endDate time.Time
	if request.StartDate != nil && len(*request.StartDate) > 0 {
		startDate, err = time.Parse(layout, *request.StartDate)
		if err != nil {
			return Task{}, ErrInvalidDate
		}
	}
	if request.EndDate != nil && len(*request.EndDate) > 0 {
		endDate, err = time.Parse(layout, *request.EndDate)
		if err != nil {
			return Task{}, ErrInvalidDate
		}
	}
	

	taskID := uuid.New()
	task := Task{
		ID: taskID,
		Description : request.Description,
		Notes: request.Notes,
		UserID: userID,
		Status: TaskStatus(request.Status),
		Priority: TaskPriority(request.Priority),
		ProjectID: projectID,
		StartDate: &startDate,
		EndDate: &endDate,
	}
	return task, nil
}