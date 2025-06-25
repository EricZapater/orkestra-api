package costitems

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CostItemRepository interface {
	Create(ctx context.Context, costItem CostItem) (CostItem, error)
	Update(ctx context.Context, costItem CostItem) (CostItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (CostItem, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]CostItem, error)
	FindAll(ctx context.Context) ([]CostItem, error)
}

type costItemRepository struct {
	db *sql.DB
}

func NewCostItemRepository(db *sql.DB) CostItemRepository {
	return &costItemRepository{
		db:db,
	}
}

func(r *costItemRepository) Create(ctx context.Context, costItem CostItem) (CostItem, error){
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cost_items(id, project_id, amount, short_description, notes, date)
		VALUES($1, $2, $3, $4, $5, $6)
	`, costItem.ID, costItem.ProjectID, costItem.Amount, costItem.ShortDescription, costItem.Notes, costItem.Date)
	if err != nil {
		return CostItem{}, err
	}
	return costItem, nil
}

func(r *costItemRepository) Update(ctx context.Context, costItem CostItem) (CostItem, error){
	_, err := r.db.ExecContext(ctx,`
	UPDATE cost_items
		SET project_id = $1,
		amount = $2,
		short_description = $3,
		notes = $4,
		date = $5
	WHERE id = $6
	`, costItem.ProjectID, costItem.Amount, costItem.ShortDescription, costItem.Notes, costItem.Date, costItem.ID)
	if err != nil {
		return CostItem{}, err
	}
	return costItem, nil
}

func(r *costItemRepository) Delete(ctx context.Context, id uuid.UUID) error{
	_, err := r.db.ExecContext(ctx,`DELETE FROM cost_items WHERE id = $1`, id)
	return err
}
func(r *costItemRepository) FindByID(ctx context.Context, id uuid.UUID) (CostItem, error){
	var costItem CostItem
	err := r.db.QueryRowContext(ctx,`
	SELECT id, project_id, amount, short_description, notes, date FROM cost_items WHERE id = $1
	`, id).Scan(&costItem.ID, &costItem.ProjectID, &costItem.Amount, &costItem.ShortDescription, &costItem.Notes, &costItem.Date)
	if err != nil {
		return CostItem{}, err
	}
	return costItem, nil

}

func(r *costItemRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]CostItem, error){
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return costItems, nil
}

func(r *costItemRepository) FindAll(ctx context.Context) ([]CostItem, error){
	var costItems []CostItem
	rows, err := r.db.QueryContext(ctx,`
	SELECT id, project_id, amount, short_description, notes, date FROM cost_items 
	`)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return costItems, nil
}