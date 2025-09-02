package projects

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, project Project) (Project, error)
	Update(ctx context.Context, project Project) (Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (Project, error)
	FindAll(ctx context.Context) ([]Project, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Project, error)
	FindBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]Project, error)
	FindCalendarBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]ProjectCalendarResponse, error)
	FindCalendarBetweenDatesByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time)([]ProjectCalendarResponse, error)
	FindProjectIDByOperatorToProjectID(ctx context.Context, operator_to_project_id uuid.UUID)(uuid.UUID, error)
	AddOperator(ctx context.Context, request OperatorToProject)([]OperatorToProject, error)
	RemoveOperator(ctx context.Context, project_id, operator_to_project_id uuid.UUID)([]OperatorToProject, error)
	FindOperatorsByProjectID(ctx context.Context, project_id uuid.UUID)([]OperatorToProject, error)
	FindOperatorsCalendarBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]ProjectCalendarResponse,error)
	AddCostItem(ctx context.Context, costItem CostItem) (CostItem, error)
	RemoveCostItem(ctx context.Context, id uuid.UUID) error
	FindCostItemsByProjectID(ctx context.Context, projectID uuid.UUID) ([]CostItem, error)
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{
		db:db,
	}
}

func(r *projectRepository) Create(ctx context.Context, project Project) (Project, error){
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO projects(id, description, start_date, end_date, color, customer_id, amount, estimated_cost)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`, project.ID, project.Description, project.StartDate, project.EndDate, project.Color, project.CustomerID, project.Amount, project.EstimatedCost)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
func(r *projectRepository) Update(ctx context.Context, project Project) (Project, error){
	_, err := r.db.ExecContext(ctx, `
		UPDATE projects
		SET description = $1,
			start_date = $2,
			end_date = $3,
			color = $4,
			customer_id = $5,
			amount = $6,
			estimated_cost =$7
		WHERE id = $8
		`, project.Description, project.StartDate, project.EndDate, project.Color, project.CustomerID, project.Amount, project.EstimatedCost, project.ID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
func(r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error{
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}
func(r *projectRepository) FindById(ctx context.Context, id uuid.UUID) (Project, error){
	var project Project
	err := r.db.QueryRowContext(ctx, `SELECT id, description, start_date, end_date, color, customer_id, amount, estimated_cost FROM projects WHERE id = $1`, 
	id,
	).Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID, &project.Amount, &project.EstimatedCost)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
func(r *projectRepository) FindAll(ctx context.Context) ([]Project, error){
	var projects []Project
	rows, err := r.db.QueryContext(ctx, `
	SELECT id, description, start_date, end_date, color, customer_id, amount, estimated_cost FROM projects
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		if err := rows.Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID, &project.Amount, &project.EstimatedCost); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func(r *projectRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Project, error){
	var projects []Project
	rows, err := r.db.QueryContext(ctx, `
	SELECT p.id, p.description, p.start_date, p.end_date, p.color, p.customer_id, p.amount, p.estimated_cost 
	FROM projects p
	INNER JOIN customer_users cu ON p.customer_id = cu.customer_id
	WHERE cu.user_id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		if err := rows.Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID, &project.Amount, &project.EstimatedCost); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func(r *projectRepository) FindBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]Project, error){
	var projects []Project
	rows, err := r.db.QueryContext(ctx, `
	SELECT id, description, start_date, end_date, color, customer_id, amount, estimated_cost FROM projects WHERE start_date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		if err := rows.Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID, &project.Amount, &project.EstimatedCost); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func(r *projectRepository) FindCalendarBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]ProjectCalendarResponse, error){
	var projects []ProjectCalendarResponse
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, 
				'[' || COALESCE(c.comercial_name, '')::text || '] ' || p.description::text  AS title, 
				p.start_date, p.end_date, p.color
				FROM projects p
				LEFT JOIN customers c ON p.customer_id = c.id
		WHERE start_date BETWEEN $1 AND $2 OR end_date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project ProjectCalendarResponse
		if err := rows.Scan(&project.ID, &project.Title, &project.StartDate, &project.EndDate, &project.Color); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}


func(r *projectRepository) FindCalendarBetweenDatesByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time)([]ProjectCalendarResponse, error){
	var projects []ProjectCalendarResponse
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, 
				'[' || COALESCE(c.comercial_name, '')::text || '] ' || p.description::text  AS title, 
				p.start_date, p.end_date, p.color
				FROM projects p
				INNER JOIN customer_users cu ON p.customer_id = cu.customer_id
				LEFT JOIN customers c ON p.customer_id = c.id
		WHERE cu.user_id = $1 AND (start_date BETWEEN $2 AND $3 OR end_date BETWEEN $2 AND $3)
	`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project ProjectCalendarResponse
		if err := rows.Scan(&project.ID, &project.Title, &project.StartDate, &project.EndDate, &project.Color); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func(r *projectRepository) FindProjectIDByOperatorToProjectID(ctx context.Context, operator_to_project_id uuid.UUID)(uuid.UUID, error){
	var project_id uuid.UUID
	err := r.db.QueryRowContext(ctx, `SELECT project_id FROM operators_to_projects WHERE id = $1`, operator_to_project_id).Scan(&project_id)
	if err != nil {
		return uuid.Nil, err
	}
	return project_id, nil
}

func(r *projectRepository) AddOperator(ctx context.Context, request OperatorToProject)([]OperatorToProject, error){
 _, err := r.db.ExecContext(ctx, `
	INSERT INTO operators_to_projects(id, operator_id, project_id, cost, dedication_percent, start_date, end_date)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	`, request.ID, request.OperatorID, request.ProjectID, request.Cost, request.DedicationPercent, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	operatorsList, err := r.FindOperatorsByProjectID(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}
	return operatorsList, nil
}
func(r *projectRepository) RemoveOperator(ctx context.Context, project_id, operator_to_project_id uuid.UUID)([]OperatorToProject, error){
	_, err := r.db.ExecContext(ctx, `DELETE FROM operators_to_projects WHERE ID`, operator_to_project_id)
	if err != nil {
		return nil, err
	}
	operatorsList, err := r.FindOperatorsByProjectID(ctx, project_id)
	if err != nil {
		return nil, err
	}
	return operatorsList, nil
}
func(r *projectRepository) FindOperatorsByProjectID(ctx context.Context, project_id uuid.UUID)([]OperatorToProject, error){
	var operatorsList []OperatorToProject
	rows, err := r.db.QueryContext(ctx, `SELECT otp.id, otp.operator_id, otp.project_id, otp.cost, otp.dedication_percent, otp.start_date, otp.end_date
	FROM operators_to_projects otp
	WHERE otp.project_id = $1`, project_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var operator OperatorToProject
		if err := rows.Scan(&operator.ID, &operator.OperatorID, &operator.ProjectID, &operator.Cost, &operator.DedicationPercent, &operator.StartDate, &operator.EndDate); err != nil{
			return nil, err
		}
		operatorsList = append(operatorsList, operator)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return operatorsList, nil
}

func(r *projectRepository) FindOperatorsCalendarBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]ProjectCalendarResponse, error){
	var projects []ProjectCalendarResponse
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, CONCAT('[', o.surname, ',', o.name,' ', op.dedication_percent,'% - ', p.description, ']')::text as title,
			op.start_date, op.end_date,
			o.color
			FROM projects p
			INNER JOIN operators_to_projects op ON p.id = op.project_id
			INNER JOIN operators o ON op.operator_id = o.id
		WHERE op.start_date BETWEEN $1 AND $2 OR op.end_date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project ProjectCalendarResponse
		if err := rows.Scan(&project.ID, &project.Title, &project.StartDate, &project.EndDate, &project.Color); err != nil{
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func(r *projectRepository) AddCostItem(ctx context.Context, costItem CostItem) (CostItem, error){
		_, err := r.db.ExecContext(ctx, `
		INSERT INTO cost_items(id, project_id, amount, short_description, notes, date)
		VALUES($1, $2, $3, $4, $5, $6)
	`, costItem.ID, costItem.ProjectID, costItem.Amount, costItem.ShortDescription, costItem.Notes, costItem.Date)
	if err != nil {
		return CostItem{}, err
	}
	return costItem, nil
}

func(r *projectRepository) RemoveCostItem(ctx context.Context, id uuid.UUID) error{
	_, err := r.db.ExecContext(ctx,`DELETE FROM cost_items WHERE id = $1`, id)
	return err
}

func(r *projectRepository) FindCostItemsByProjectID(ctx context.Context, projectID uuid.UUID) ([]CostItem, error){
	var costItems []CostItem
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, project_id, amount, short_description, notes, date FROM cost_items WHERE project_id = $1
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var costItem CostItem
		if err := rows.Scan(&costItem.ID, &costItem.ProjectID, &costItem.Amount, &costItem.ShortDescription, &costItem.Notes, &costItem.Date); err != nil {
			return nil, err
		}
		costItems = append(costItems, costItem)
	}
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return costItems, nil
}