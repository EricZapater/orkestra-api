package menus

import (
	"context"

	"github.com/google/uuid"
)

type MenuService interface {
	Create(ctx context.Context, request *CreateMenuRequest) (Menu, error)
	Update(ctx context.Context, id string, request *CreateMenuRequest) (Menu, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Menu, error)
	FindByProfileID(ctx context.Context, id string) ([]Menu, error)
	FindAll(ctx context.Context) ([]Menu, error)
}

type menuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) Create(ctx context.Context, request *CreateMenuRequest) (Menu, error) {
	parentUUID, err := uuid.Parse(request.ParentID)
	if err != nil {
		return Menu{}, ErrInvalidParentID
	}
	menu := &Menu{
		ID:          uuid.New(),
		Label:       request.Label,
		Icon:        request.Icon,
		Route:       request.Route,
		ParentID:    parentUUID,
		SortOrder:   request.SortOrder,
		IsSeparator: request.IsSeparator,
	}

	return s.repo.Create(ctx, menu)
}

func (s *menuService) Update(ctx context.Context, id string, request *CreateMenuRequest) (Menu, error) {
	parentUUID, err := uuid.Parse(request.ParentID)
	if err != nil {
		return Menu{}, ErrInvalidParentID
	}
	menuUUID, err := uuid.Parse(id)
	if err != nil {
		return Menu{}, ErrInvalidMenuID
	}
	menu := &Menu{
		ID:          menuUUID,
		Label:       request.Label,
		Icon:        request.Icon,
		Route:       request.Route,
		ParentID:    parentUUID,
		SortOrder:   request.SortOrder,
		IsSeparator: request.IsSeparator,
	}

	return s.repo.Update(ctx, menu)
}

func (s *menuService) Delete(ctx context.Context, id string) error {
	menuUUID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidMenuID
	}
	return s.repo.Delete(ctx, menuUUID)
}

func (s *menuService) FindByID(ctx context.Context, id string) (Menu, error) {
	menuUUID, err := uuid.Parse(id)
	if err != nil {
		return Menu{}, ErrInvalidMenuID
	}
	return s.repo.FindByID(ctx, menuUUID)
}

func (s *menuService) FindByProfileID(ctx context.Context, id string) ([]Menu, error) {
	profileUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidProfileID
	}
	return s.repo.FindByProfileID(ctx, profileUUID)
}

func (s *menuService) FindAll(ctx context.Context) ([]Menu, error) {
	return s.repo.FindAll(ctx)
}