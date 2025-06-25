package menus

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type MenuRepository interface {
	Create(ctx context.Context, request *Menu)(Menu, error)
	Update(ctx context.Context, request *Menu)(Menu, error)
	Delete(ctx context.Context, id uuid.UUID)error
	FindByID(ctx context.Context, id uuid.UUID)(Menu, error)
	FindByProfileID(ctx context.Context, id uuid.UUID)([]Menu, error)
	FindAll(ctx context.Context)([]Menu, error)
}

type menuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) MenuRepository{
	return &menuRepository{db}
}

func(r *menuRepository) Create(ctx context.Context, request *Menu)(Menu, error){
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO menu_items(ID, label, icon, route, parent_id, sort_order, is_separator)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	`, request.ID, request.Label, request.Icon, request.Route, request.ParentID, request.SortOrder, request.IsSeparator)

	if err != nil {
		return Menu{}, err
	}
	return *request, nil
}
func(r *menuRepository) Update(ctx context.Context, request *Menu)(Menu, error){
	_, err := r.db.ExecContext(ctx, `
	UPDATE menu_items
		SET label = $1,
		icon = $2,
		route = $3,
		parent_id = $4,
		sort_order = $5,
		is_separator = $6	
	WHERE ID = $7
	`, request.Label, request.Icon, request.Route, request.ParentID, request.SortOrder, request.IsSeparator, request.ID)
	if err != nil {
		return Menu{}, err
	}
	return *request, nil
}
func(r *menuRepository) Delete(ctx context.Context, id uuid.UUID)error{
	_, err := r.db.ExecContext(ctx, `
	DELETE FROM menu_items WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
func(r *menuRepository) FindByID(ctx context.Context, id uuid.UUID)(Menu, error){
	var menu Menu
	query := `SELECT ID, label, icon, route, parent_id, sort_order, is_separator FROM menu_items WHERE id = $1 `
	err := r.db.QueryRowContext(ctx, query, id).Scan(&menu.ID, &menu.Label, &menu.Icon, &menu.Route, &menu.ParentID, &menu.SortOrder, &menu.IsSeparator)
	if err != nil {
		return Menu{}, nil
	}
	return menu, nil
}
func(r *menuRepository) FindByProfileID(ctx context.Context, profileID uuid.UUID)([]Menu, error){
	var menus []Menu
	query := `SELECT m.ID, m.label, m.icon, m.route, m.parent_id, m.sort_order, m.is_separator 
	FROM menu_items m
		INNER JOIN profile_menus p ON p.menu_id = m.id
	WHERE p.profile_id = $1
	ORDER BY m.sort_order
	`
	rows, err := r.db.QueryContext(ctx, query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.Label, &menu.Icon, &menu.Route, &menu.ParentID, &menu.SortOrder, &menu.IsSeparator); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}
func(r *menuRepository) FindAll(ctx context.Context)([]Menu, error){
	var menus []Menu
	query := `SELECT m.ID, m.label, m.icon, m.route, m.parent_id, m.sort_order, m.is_separator 
	FROM menu_items m
	ORDER BY m.sort_order`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.Label, &menu.Icon, &menu.Route, &menu.ParentID, &menu.SortOrder, &menu.IsSeparator); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}