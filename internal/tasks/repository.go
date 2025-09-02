package tasks

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task Task) (Task, error)
	Update(ctx context.Context, task Task) (Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (Task, error)
	FindByStatus(ctx context.Context, status TaskStatus) ([]Task, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Task, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]Task, error)
	FindByPriority(ctx context.Context, priority TaskPriority) ([]Task, error)
	FindAll(ctx context.Context) ([]Task, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository{
	return &taskRepository{
		db : db,
	}
}

func (r *taskRepository) Create(ctx context.Context, task Task) (Task, error) {
	_, err := r.db.ExecContext(ctx,`
	INSERT INTO tasks(id, description, notes, user_id, status, priority, project_id, start_date, end_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, task.ID, task.Description, task.Notes, task.UserID, task.Status, task.Priority, task.ProjectID, task.StartDate, task.EndDate)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}
func (r *taskRepository) Update(ctx context.Context, task Task) (Task, error){
	_, err := r.db.ExecContext(ctx,`
	UPDATE  tasks
	SET description = $1, 
		notes = $2, 
		user_id = $3, 
		status = $4,  
		priority = $5, 
		project_id = $6,  
		start_date = $7, 
		end_date = $8
	WHERE ID = $9
	`, task.Description, task.Notes, task.UserID, task.Status, task.Priority, task.ProjectID, task.StartDate, task.EndDate, task.ID)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}
func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error{
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}
func (r *taskRepository) FindById(ctx context.Context, id uuid.UUID) (Task, error){
	var task Task
	err := r.db.QueryRowContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks WHERE id = $1
	`, id,
	).Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *taskRepository) FindByStatus(ctx context.Context, status TaskStatus) ([]Task, error){
	var tasks []Task
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks WHERE status = $1
	`, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *taskRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Task, error){
	var tasks []Task
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks WHERE user_id = $1
	`, userID,
)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *taskRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]Task, error){
	var tasks []Task
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks WHERE project_id = $1
	`, projectID,
)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *taskRepository) FindByPriority(ctx context.Context, priority TaskPriority) ([]Task, error){
	var tasks []Task
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks WHERE priority = $1
	`, priority,
)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *taskRepository) FindAll(ctx context.Context) ([]Task, error){
	var tasks []Task
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, description, notes, user_id, status, priority, project_id, start_date, end_date
	FROM tasks ORDER BY project_id, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Notes, &task.UserID, &task.Status, &task.Priority, &task.ProjectID, &task.StartDate, &task.EndDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}