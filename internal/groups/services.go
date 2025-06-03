package groups

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GroupService interface {
	Create(ctx context.Context, request GroupRequest) (Group, error)
	Update(ctx context.Context, id string, request GroupRequest) (Group, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Group, error)
	FindByName(ctx context.Context, name string) (Group, error)
	FindAll(ctx context.Context) ([]Group, error)
	FindByUserID(ctx context.Context, userID string) ([]Group, error)
	AddUserToGroup(ctx context.Context, groupID, userID string) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID string) error
}

type groupService struct {
	repo GroupRepository
}

func NewGroupService(repo GroupRepository) GroupService {
	return &groupService{repo}
}

func (s *groupService) Create(ctx context.Context, request GroupRequest) (Group, error) {
	// Validate the request
	if request.Name == "" {
		return Group{}, ErrInvalidRequest
	}

	// Check if the group name is already taken
	_, err := s.repo.FindByName(ctx, request.Name)
	if err == nil {
		return Group{}, ErrGroupNameTaken
	}
	if err != ErrGroupNotFound {
		return Group{}, err
	}

	group := Group{
		ID:   uuid.New(),
		Name: request.Name,
	}

	return s.repo.Create(ctx, group)
}

func (s *groupService) Update(ctx context.Context, id string, request GroupRequest) (Group, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return Group{}, ErrInvalidID
	}

	// Validate the request
	if request.Name == "" {
		return Group{}, ErrInvalidRequest
	}

	group, err := s.repo.FindByID(ctx, groupID)
	if err != nil {
		return Group{}, err
	}

	group.Name = request.Name

	return s.repo.Update(ctx, group)
}

func (s *groupService) Delete(ctx context.Context, id string) error {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	return s.repo.Delete(ctx, groupID)
}

func (s *groupService) FindByID(ctx context.Context, id string) (Group, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return Group{}, ErrInvalidID
	}

	group, err := s.repo.FindByID(ctx, groupID)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (s *groupService) FindByName(ctx context.Context, name string) (Group, error) {
	group, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (s *groupService) FindAll(ctx context.Context) ([]Group, error) {
	groups, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *groupService) FindByUserID(ctx context.Context, userID string) ([]Group, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidID
	}

	groups, err := s.repo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *groupService) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	fmt.Println("AddUserToGroup")
	fmt.Printf("groupID: %s, userID: %s\n", groupID, userID)
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return ErrInvalidID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidID
	}
	id := uuid.New()
	return s.repo.AddUserToGroup(ctx, id, groupUUID, userUUID)
}

func (s *groupService) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return ErrInvalidID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidID
	}

	return s.repo.RemoveUserFromGroup(ctx, groupUUID, userUUID)
}