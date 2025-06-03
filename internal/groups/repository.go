package groups

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type GroupRepository interface {
	Create(ctx context.Context, group Group) (Group, error)
	Update(ctx context.Context, group Group) (Group, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (Group, error)
	FindByName(ctx context.Context, name string) (Group, error)
	FindAll(ctx context.Context) ([]Group, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Group, error)
	AddUserToGroup(ctx context.Context, id, groupID, userID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error
}

type groupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) GroupRepository {
	return &groupRepository{db}
}

func (r *groupRepository) Create(ctx context.Context, group Group) (Group, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO groups (id, name, created_at)
		VALUES ($1, $2, NOW())`,
		group.ID, group.Name,
	)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func (r *groupRepository) Update(ctx context.Context, group Group) (Group, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE groups
		SET name = $1
		WHERE id = $2`,
		group.Name, group.ID,
	)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func (r *groupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM groups
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *groupRepository) FindByID(ctx context.Context, id uuid.UUID) (Group, error) {
	var group Group
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, created_at
		FROM groups
		WHERE id = $1`,
		id,
	).Scan(&group.ID, &group.Name, &group.CreatedAt)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func (r *groupRepository) FindByName(ctx context.Context, name string) (Group, error) {
	var group Group
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, created_at
		FROM groups
		WHERE name = $1`,
		name,
	).Scan(&group.ID, &group.Name, &group.CreatedAt)
	if err == sql.ErrNoRows {
		return Group{}, ErrGroupNotFound
	}else if err != nil {
		return Group{}, err
	}
	return group, nil
}

func (r *groupRepository) FindAll(ctx context.Context) ([]Group, error) {
	var groups []Group
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at
		FROM groups`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Group, error) {
	var groups []Group
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id, g.name, g.created_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepository) AddUserToGroup(ctx context.Context,id, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO group_members (id, group_id, user_id, joined_at)
		VALUES ($1, $2, $3, NOW())`,
		id, groupID, userID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *groupRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM group_members
		WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	)
	if err != nil {
		return err
	}
	return nil
}