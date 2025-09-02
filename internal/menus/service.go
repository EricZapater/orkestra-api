package menus

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type MenuService interface {
	Create(ctx context.Context, request *CreateMenuRequest) (Menu, error)
	Update(ctx context.Context, id string, request *CreateMenuRequest) (Menu, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Menu, error)
	FindByProfileID(ctx context.Context, id string) ([]MenuTree, error)
	FindAll(ctx context.Context) ([]Menu, error)
	AddMenuToProfile(ctx context.Context, request *MenuToProfileRequest) error
	RemoveMenuFromProfile(ctx context.Context, request *MenuToProfileRequest) error
}

type menuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) Create(ctx context.Context, request *CreateMenuRequest) (Menu, error) {
	var parentUUID *uuid.UUID // ← Pointer, per defecte nil
    
    if request.ParentID != nil && *request.ParentID != "" {
        parsed, err := uuid.Parse(*request.ParentID)
        if err != nil {
            return Menu{}, ErrInvalidParentID
        }
        parentUUID = &parsed // ← Assigna l'adreça del UUID parsed
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
	fmt.Println(menu)
	return s.repo.Create(ctx, menu)
}

func (s *menuService) Update(ctx context.Context, id string, request *CreateMenuRequest) (Menu, error) {
	var parentUUID *uuid.UUID // ← Pointer, per defecte nil
    
    if request.ParentID != nil && *request.ParentID != "" {
        parsed, err := uuid.Parse(*request.ParentID)
        if err != nil {
            return Menu{}, ErrInvalidParentID
        }
        parentUUID = &parsed // ← Assigna l'adreça del UUID parsed
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

func (s *menuService) FindByProfileID(ctx context.Context, id string) ([]MenuTree, error) {
	profileUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidProfileID
	}
	menus, err := s.repo.FindByProfileID(ctx, profileUUID)
	if err !=nil {
		return nil, err
	}

	menutree := buildMenuTree(menus)
	sortChildrenRecursive(menutree)
	return menutree, nil
}

func (s *menuService) FindAll(ctx context.Context) ([]Menu, error) {
	return s.repo.FindAll(ctx)
}

func (s *menuService) AddMenuToProfile(ctx context.Context, request *MenuToProfileRequest) error {
	profileUUID, err := uuid.Parse(request.ProfileID)
	if err != nil {
		return ErrInvalidProfileID
	}
	menuUUID, err := uuid.Parse(request.MenuID)
	if err != nil {
		return ErrInvalidMenuID
	}
	return s.repo.AddMenuToProfile(ctx, profileUUID, menuUUID)
}

func (s *menuService) RemoveMenuFromProfile(ctx context.Context, request *MenuToProfileRequest) error {
	profileUUID, err := uuid.Parse(request.ProfileID)
	if err != nil {
		return ErrInvalidProfileID
	}
	menuUUID, err := uuid.Parse(request.MenuID)
	if err != nil {
		return ErrInvalidMenuID
	}
	return s.repo.RemoveMenuFromProfile(ctx, profileUUID, menuUUID)
}

func buildMenuTree(menus []Menu) []MenuTree {
	// Map per accedir ràpidament per ID
	menuMap := make(map[uuid.UUID]*MenuTree)
	var roots []MenuTree

	// Inicialitzem els nodes
	for _, m := range menus {
		node := MenuTree{
			ID:          m.ID,
			Label:       m.Label,
			Icon:        m.Icon,
			Route:       m.Route,
			ParentID:    m.ParentID,
			SortOrder:   m.SortOrder,
			IsSeparator: m.IsSeparator,
			Children:    []MenuTree{},
		}
		menuMap[m.ID] = &node
	}
	
	for _, node := range menuMap {
		if node.ParentID == nil {
			roots = append(roots, *node)
		} else if parent, ok := menuMap[*node.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}
	
	for i := range roots {		
		sort.SliceStable(roots[i].Children, func(a, b int) bool {
			return roots[i].Children[a].SortOrder < roots[i].Children[b].SortOrder
		})
	}

	return roots
}

func sortChildrenRecursive(nodes []MenuTree) {
    for i := range nodes {
        sort.SliceStable(nodes[i].Children, func(a, b int) bool {
            return nodes[i].Children[a].SortOrder < nodes[i].Children[b].SortOrder
        })
        sortChildrenRecursive(nodes[i].Children)
    }
}