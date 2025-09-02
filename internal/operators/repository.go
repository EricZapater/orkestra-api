package operators

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type OperatorRepository interface {
	Create(ctx context.Context, operator Operator)(Operator, error)
	Update(ctx context.Context, operator Operator)(Operator, error)
	Delete(ctx context.Context, id uuid.UUID)error
	FindByID(ctx context.Context, id uuid.UUID)(Operator, error)
	FindAll(ctx context.Context)([]Operator, error)
}

type operatorRepository struct {
	db *sql.DB
}

func NewOperatorRepository(db *sql.DB) OperatorRepository{
	return &operatorRepository{
		db:db,
	}
}

func(r *operatorRepository) Create(ctx context.Context, operator Operator)(Operator, error){
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO operators(id, name, surname, cost, color)
	VALUES($1, $2, $3, $4, $5)
	`, operator.ID, operator.Name, operator.Surname, operator.Cost, operator.Color)
	if err != nil {
		return Operator{}, err
	}
	return operator, nil
}

func(r *operatorRepository) Update(ctx context.Context, operator Operator)(Operator, error){
	_, err := r.db.ExecContext(ctx, `
	UPDATE operators
	SET name = $1,
		surname = $2,
		cost = $3,
		color = $4
	WHERE id = $5
	`, operator.Name, operator.Surname, operator.Cost, operator.Color, operator.ID)
	if err != nil {
		return Operator{}, nil
	}
	return operator, nil
}

func(r *operatorRepository) Delete(ctx context.Context, id uuid.UUID)error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM operators WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func(r *operatorRepository) FindByID(ctx context.Context, id uuid.UUID)(Operator, error){
	var operator Operator
	err := r.db.QueryRowContext(ctx, 
		`SELECT ID, name, surname, cost, color FROM operators WHERE ID = $1`,
		id).Scan(&operator.ID, &operator.Name, &operator.Surname, &operator.Cost, &operator.Color)

	if err != nil {
		return Operator{}, err
	}
	return operator, nil
}

func(r *operatorRepository) FindAll(ctx context.Context)([]Operator, error){
	var operators []Operator
	rows, err := r.db.QueryContext(ctx, `SELECT ID, name, surname, cost, color FROM operators`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var operator Operator
		if err := rows.Scan(&operator.ID, &operator.Name, &operator.Surname, &operator.Cost, &operator.Color); err != nil{
			return nil, err
		}
		operators = append(operators, operator)
	}
	return operators, nil
}