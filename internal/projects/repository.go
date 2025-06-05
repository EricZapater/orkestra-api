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
	FindBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]Project, error)
	FindCalendarBetweenDates(ctx context.Context, startDate, endDate *time.Time)([]ProjectCalendarResponse, error)
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
	INSERT INTO projects(id, description, start_date, end_date, color, customer_id)
	VALUES($1, $2, $3, $4, $5, $6)
	`, project.ID, project.Description, project.StartDate, project.EndDate, project.Color, project.CustomerID)
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
			customer_id = $5
		WHERE id = $6
		`, project.Description, project.StartDate, project.EndDate, project.Color, project.CustomerID, project.ID)
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
	err := r.db.QueryRowContext(ctx, `SELECT id, description, start_date, end_date, color, customer_id FROM projects WHERE id = $1`, 
	id,
	).Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
func(r *projectRepository) FindAll(ctx context.Context) ([]Project, error){
	var projects []Project
	rows, err := r.db.QueryContext(ctx, `
	SELECT id, description, start_date, end_date, color, customer_id FROM projects
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		if err := rows.Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID); err != nil{
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
	SELECT id, description, start_date, end_date, color, customer_id FROM projects WHERE start_date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		if err := rows.Scan(&project.ID, &project.Description, &project.StartDate, &project.EndDate, &project.Color, &project.CustomerID); err != nil{
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
		WHERE start_date BETWEEN $1 AND $2
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