package category

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	categorydto "sys-admin-serve/internal/dto/category"
	"sys-admin-serve/internal/model"
	repositorycategory "sys-admin-serve/internal/repository/category"

	"go.uber.org/zap"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrCategoryInvalidRequest = errors.New("invalid category request")
	ErrCategoryNotFound       = errors.New("category not found")
	ErrCategoryParentNotFound = errors.New("parent category not found")
	ErrCategoryInvalidParent  = errors.New("invalid parent category")
	ErrCategoryCodeExists     = errors.New("category code already exists")
	ErrCategoryNameExists     = errors.New("category name already exists")
)

type Service struct {
	repo *repositorycategory.Repository
	log  *zap.Logger
}

func NewService(repo *repositorycategory.Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) ListCategories(ctx context.Context, req categorydto.ListCategoriesRequest) (*categorydto.ListCategoriesResult, error) {
	page := req.Page
	if page <= 0 {
		page = defaultPage
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	categories, total, err := s.repo.ListCategories(ctx, repositorycategory.ListFilter{
		Keyword:  strings.TrimSpace(req.Keyword),
		Status:   req.Status,
		ParentID: req.ParentID,
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	result := &categorydto.ListCategoriesResult{
		List:     make([]categorydto.CategoryItem, 0, len(categories)),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	for i := range categories {
		result.List = append(result.List, convertCategoryItem(categories[i]))
	}

	return result, nil
}

func (s *Service) GetCategoryTree(ctx context.Context, req categorydto.CategoryTreeRequest) ([]categorydto.CategoryTreeItem, error) {
	parentID := uint64(0)
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	if parentID != 0 {
		parent, err := s.repo.GetCategoryByID(ctx, parentID)
		if err != nil {
			return nil, fmt.Errorf("get parent category by id: %w", err)
		}
		if parent == nil {
			return nil, ErrCategoryParentNotFound
		}
	}

	categories, err := s.repo.ListCategoriesForTree(ctx, repositorycategory.TreeFilter{
		Status: req.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("list categories for tree: %w", err)
	}

	nodeMap := make(map[uint64]*categorydto.CategoryTreeItem, len(categories))
	roots := make([]*categorydto.CategoryTreeItem, 0)
	for i := range categories {
		category := categories[i]
		nodeMap[category.ID] = &categorydto.CategoryTreeItem{
			ID:       category.ID,
			ParentID: category.ParentID,
			Name:     category.Name,
			Code:     category.Code,
			Sort:     category.Sort,
			Status:   category.Status,
			Icon:     category.Icon,
			Remark:   category.Remark,
			Children: make([]categorydto.CategoryTreeItem, 0),
		}
	}

	for i := range categories {
		node := nodeMap[categories[i].ID]
		if node.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := nodeMap[node.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, *node)
	}

	for i := range roots {
		sortCategoryTreeChildren(&roots[i].Children)
	}
	sort.SliceStable(roots, func(i, j int) bool {
		if roots[i].Sort != roots[j].Sort {
			return roots[i].Sort < roots[j].Sort
		}
		return roots[i].ID < roots[j].ID
	})

	if parentID == 0 {
		return dereferenceTreeNodes(roots), nil
	}

	parent, ok := nodeMap[parentID]
	if !ok {
		return []categorydto.CategoryTreeItem{}, nil
	}

	sortCategoryTreeChildren(&parent.Children)
	return parent.Children, nil
}

func (s *Service) CreateCategory(ctx context.Context, req categorydto.CreateCategoryRequest) (*categorydto.CategoryItem, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" || code == "" {
		return nil, ErrCategoryInvalidRequest
	}

	if err := s.ensureParentValid(ctx, req.ParentID, 0); err != nil {
		return nil, err
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	category := &model.Category{
		ParentID: req.ParentID,
		Name:     name,
		Code:     code,
		Sort:     req.Sort,
		Status:   status,
		Icon:     strings.TrimSpace(req.Icon),
		Remark:   strings.TrimSpace(req.Remark),
	}
	if err := s.repo.CreateCategory(ctx, category); err != nil {
		switch {
		case errors.Is(err, repositorycategory.ErrDuplicateCategoryCode):
			return nil, ErrCategoryCodeExists
		case errors.Is(err, repositorycategory.ErrDuplicateCategoryName):
			return nil, ErrCategoryNameExists
		default:
			return nil, fmt.Errorf("create category: %w", err)
		}
	}

	s.log.Info("category created", zap.Uint64("category_id", category.ID), zap.String("category_code", category.Code))
	item := convertCategoryItem(*category)
	return &item, nil
}

func (s *Service) UpdateCategory(ctx context.Context, categoryID uint64, req categorydto.UpdateCategoryRequest) (*categorydto.CategoryItem, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" || code == "" {
		return nil, ErrCategoryInvalidRequest
	}

	existing, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	if existing == nil {
		return nil, ErrCategoryNotFound
	}

	if err := s.ensureParentValid(ctx, req.ParentID, categoryID); err != nil {
		return nil, err
	}

	update := &model.Category{
		ID:       categoryID,
		ParentID: req.ParentID,
		Name:     name,
		Code:     code,
		Sort:     req.Sort,
		Status:   req.Status,
		Icon:     strings.TrimSpace(req.Icon),
		Remark:   strings.TrimSpace(req.Remark),
	}
	if err := s.repo.UpdateCategory(ctx, update); err != nil {
		switch {
		case errors.Is(err, repositorycategory.ErrDuplicateCategoryCode):
			return nil, ErrCategoryCodeExists
		case errors.Is(err, repositorycategory.ErrDuplicateCategoryName):
			return nil, ErrCategoryNameExists
		default:
			return nil, fmt.Errorf("update category: %w", err)
		}
	}

	s.log.Info("category updated", zap.Uint64("category_id", categoryID), zap.String("category_code", update.Code))
	item := convertCategoryItem(*update)
	return &item, nil
}

func (s *Service) ensureParentValid(ctx context.Context, parentID uint64, selfID uint64) error {
	if parentID == 0 {
		return nil
	}
	if selfID > 0 && parentID == selfID {
		return ErrCategoryInvalidParent
	}

	parent, err := s.repo.GetCategoryByID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("get parent category by id: %w", err)
	}
	if parent == nil {
		return ErrCategoryParentNotFound
	}

	return nil
}

func convertCategoryItem(category model.Category) categorydto.CategoryItem {
	return categorydto.CategoryItem{
		ID:       category.ID,
		ParentID: category.ParentID,
		Name:     category.Name,
		Code:     category.Code,
		Sort:     category.Sort,
		Status:   category.Status,
		Icon:     category.Icon,
		Remark:   category.Remark,
	}
}

func sortCategoryTreeChildren(children *[]categorydto.CategoryTreeItem) {
	sort.SliceStable(*children, func(i, j int) bool {
		if (*children)[i].Sort != (*children)[j].Sort {
			return (*children)[i].Sort < (*children)[j].Sort
		}
		return (*children)[i].ID < (*children)[j].ID
	})

	for i := range *children {
		sortCategoryTreeChildren(&(*children)[i].Children)
	}
}

func dereferenceTreeNodes(nodes []*categorydto.CategoryTreeItem) []categorydto.CategoryTreeItem {
	result := make([]categorydto.CategoryTreeItem, 0, len(nodes))
	for i := range nodes {
		result = append(result, *nodes[i])
	}

	return result
}
