package category

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"sys-admin-serve/internal/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateCategoryCode = errors.New("duplicate category code")
	ErrDuplicateCategoryName = errors.New("duplicate category name")
)

type ListFilter struct {
	Keyword  string
	Status   *int
	ParentID *uint64
	Offset   int
	Limit    int
}

type TreeFilter struct {
	Status *int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCategoryByID(ctx context.Context, categoryID uint64) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", categoryID).
		First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query category by id: %w", err)
	}

	return &category, nil
}

func (r *Repository) ListCategories(ctx context.Context, filter ListFilter) ([]model.Category, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Category{}).Where("deleted_at IS NULL")
	if filter.Keyword != "" {
		likeKeyword := "%" + filter.Keyword + "%"
		query = query.Where("(name LIKE ? OR code LIKE ?)", likeKeyword, likeKeyword)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count categories: %w", err)
	}

	categories := make([]model.Category, 0)
	if err := query.Order("sort ASC, id ASC").Offset(filter.Offset).Limit(filter.Limit).Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("list categories: %w", err)
	}

	return categories, total, nil
}

func (r *Repository) CreateCategory(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return mapCategoryWriteError(err, "create category")
	}

	return nil
}

func (r *Repository) ListCategoriesForTree(ctx context.Context, filter TreeFilter) ([]model.Category, error) {
	query := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("deleted_at IS NULL")
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	categories := make([]model.Category, 0)
	if err := query.Order("parent_id ASC, sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("list categories for tree: %w", err)
	}

	return categories, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, category *model.Category) error {
	updates := map[string]any{
		"parent_id": category.ParentID,
		"name":      category.Name,
		"code":      category.Code,
		"sort":      category.Sort,
		"status":    category.Status,
		"icon":      category.Icon,
		"remark":    category.Remark,
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ? AND deleted_at IS NULL", category.ID).
		Updates(updates).Error; err != nil {
		return mapCategoryWriteError(err, "update category")
	}

	return nil
}

func mapCategoryWriteError(err error, action string) error {
	if !isDuplicateEntryError(err) {
		return fmt.Errorf("%s: %w", action, err)
	}

	if isDuplicateKey(err, "uk_categories_code") {
		return ErrDuplicateCategoryCode
	}
	if isDuplicateKey(err, "uk_categories_parent_name") {
		return ErrDuplicateCategoryName
	}

	return fmt.Errorf("%s: %w", action, err)
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func isDuplicateKey(err error, key string) bool {
	var mysqlErr *mysqlDriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}

	return strings.Contains(strings.ToLower(mysqlErr.Message), strings.ToLower(key))
}
