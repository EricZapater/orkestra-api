package projects

import (
	"context"
	"fmt"
	"log"
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
	AddOperator(ctx context.Context, request OperatorToProjectRequest)([]OperatorToProject, error)
	RemoveOperator(ctx context.Context, id string)([]OperatorToProject, error)
	FindOperatorsByProjectID(ctx context.Context, projectID string)([]OperatorToProject, error)
	FindOperatorsCalendarBetweenDates(ctx context.Context, startDate, endDate string)([]ProjectCalendarResponse, error)
	AddCostItem(ctx context.Context, request *CostItemRequest) (CostItem, error)
	RemoveCostItem(ctx context.Context, id string) error
	FindCostItemsByProjectID(ctx context.Context, projectID string) ([]CostItem, error)
	CalculateProjectCost(ctx context.Context, projectID string) (decimal.Decimal, error)		
	UpdateProjectCost(ctx context.Context, projectID string, newCost decimal.Decimal) error
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

func(s *projectService) AddOperator(ctx context.Context, request OperatorToProjectRequest)([]OperatorToProject, error){
	log.Default().Println("AddOperator called with request:", request)
	if request.OperatorID == "" || request.ProjectID == "" {
		return nil, ErrInvalidRequest
	}
	operatorUUID, err := uuid.Parse(request.OperatorID)
	if err != nil {
		return nil, ErrInvalidID
	}
	projectUUID, err := uuid.Parse(request.ProjectID)
	if err != nil {
		return nil, ErrInvalidID
	}
	costDecimal, err := decimal.NewFromString(request.Cost)
	if err != nil {
		return nil, ErrInvalidRequest
	}

	dedicationPercentDecimal, err := decimal.NewFromString(request.DedicationPercent)
	if err != nil {
		return nil, ErrInvalidRequest
	}
	layout := time.RFC3339 
	startDate, err := time.Parse(layout, request.StartDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	endDate, err := time.Parse(layout, request.EndDate)
	if err != nil {
		return nil, ErrInvalidDate
	}

	operatorToProject := OperatorToProject{
		ID: uuid.New(),
		OperatorID: operatorUUID,
		ProjectID: projectUUID,
		Cost: costDecimal,
		DedicationPercent: dedicationPercentDecimal,
		StartDate: startDate,
		EndDate: endDate,
	}
	operators, err := s.repo.AddOperator(ctx, operatorToProject)
	if err != nil {
		return nil, err
	}
	newCost, err := s.CalculateProjectCost(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}
	err = s.UpdateProjectCost(ctx, request.ProjectID, newCost)
	if err != nil {
		return nil, err
	}

	return operators, nil
}

func(s *projectService) RemoveOperator(ctx context.Context, id string)([]OperatorToProject, error){
	operatorUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	//Recuperar el project_id
	project_id, err := s.repo.FindProjectIDByOperatorToProjectID(ctx, operatorUUID)
	if err != nil {
		return nil, err
	}
	operators, err := s.repo.RemoveOperator(ctx,project_id,  operatorUUID)
	if err != nil {
		return nil, err
	}
	newCost, err := s.CalculateProjectCost(ctx, project_id.String())
	if err != nil {
		return nil, err
	}
	err = s.UpdateProjectCost(ctx, project_id.String(), newCost)
	if err != nil {
		return nil, err
	}

	return operators, nil
}

func(s *projectService) FindOperatorsByProjectID(ctx context.Context, projectID string)([]OperatorToProject, error){
	if projectID == "" {
		return nil, ErrInvalidRequest
	}
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindOperatorsByProjectID(ctx, projectUUID)
}

func(s *projectService)FindOperatorsCalendarBetweenDates(ctx context.Context, startDate, endDate string)([]ProjectCalendarResponse, error){
	layout := "2006-01-02" 

	pStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	pEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	return s.repo.FindOperatorsCalendarBetweenDates(ctx, &pStartDate, &pEndDate)
}

func(s *projectService) AddCostItem(ctx context.Context, request *CostItemRequest) (CostItem, error){
		_, err := s.FindById(ctx, request.ProjectID)
	if err != nil {
		return CostItem{}, ErrProjectNotFound
	}

	costItem,err := createCostItemModelFromRequest(request)
	if err != nil{
		return CostItem{}, err
	}
	costItem, err = s.repo.AddCostItem(ctx, costItem)
	if err != nil {
		return CostItem{}, err
	}
	newCost, err := s.CalculateProjectCost(ctx, request.ProjectID)
	if err != nil {
		return CostItem{}, err
	}
	err = s.UpdateProjectCost(ctx, request.ProjectID, newCost)
	if err != nil {
		return CostItem{}, err
	}
	return costItem, nil
}

func(s *projectService) RemoveCostItem(ctx context.Context, id string) error{
	costItemID, err := uuid.Parse(id)
	if err != nil{
		return ErrInvalidID
	}
	err = s.repo.Delete(ctx, costItemID)
	if err != nil {
		return  err
	}
	newCost, err := s.CalculateProjectCost(ctx, id)
	if err != nil {
		return  err
	}
	err = s.UpdateProjectCost(ctx, id, newCost)
	if err != nil {
		return err
	}
	return nil
}

func(s *projectService) FindCostItemsByProjectID(ctx context.Context, projectID string) ([]CostItem, error){
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindCostItemsByProjectID(ctx, projectUUID)
}

func (s *projectService) CalculateProjectCost(ctx context.Context, projectID string) (decimal.Decimal, error) {
	
	operators, err := s.FindOperatorsByProjectID(ctx, projectID)
	if err != nil {
		return decimal.Zero, err
	}
	totalCost := decimal.Zero
	
	for _, operator := range operators {
		workingDays := decimal.NewFromInt(int64(WorkingDaysBetween(operator.StartDate, operator.EndDate)))
		operatorCost := operator.Cost.
							Mul(workingDays).
							Mul(operator.DedicationPercent.Div(decimal.NewFromInt(100)))

		totalCost = totalCost.Add(operatorCost)
	}

	costItems, err := s.FindCostItemsByProjectID(ctx, projectID)
	if err != nil {
		return decimal.Zero, err
	}
	for _, item := range costItems {
		totalCost = totalCost.Add(item.Amount)
	}
	
	return totalCost, nil
}

func (s *projectService) UpdateProjectCost(ctx context.Context, projectID string, newCost decimal.Decimal) error {
	project, err := s.FindById(ctx, projectID)
	if err != nil {
		return ErrInvalidID
	}
	
	projectRequest := ProjectRequest{
		Color: project.Color,
		CustomerID: project.CustomerID.String(),
		Description: project.Description,
		StartDate: project.StartDate.Format(time.RFC3339),
		EndDate: project.EndDate.Format(time.RFC3339),
		Amount: project.Amount.String(),
		EstimatedCost: newCost.String(),
	}
	_, err = s.Update(ctx, project.ID.String(), projectRequest)
	if err != nil {
		return err
	}
	return nil
}


func WorkingDaysBetween(start, end time.Time) int {	
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			count++
		}
	}
	return count
}

func createCostItemModelFromRequest(request *CostItemRequest)(CostItem, error){
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