package auth

import (
	"context"
	"errors"
	"fmt"

	"sys-admin-serve/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) ListRoleCodesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var roleCodes []string
	err := r.db.WithContext(ctx).
		Table("roles").
		Select("roles.code").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.deleted_at IS NULL").
		Order("roles.sort ASC, roles.id ASC").
		Find(&roleCodes).Error
	if err != nil {
		return nil, fmt.Errorf("list role codes by user id: %w", err)
	}

	return roleCodes, nil
}

func (r *Repository) UpdateUserLastLogin(ctx context.Context, userID uint64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("CURRENT_TIMESTAMP(3)")).Error; err != nil {
		return fmt.Errorf("update user last login: %w", err)
	}

	return nil
}

func (r *Repository) ListMenusByUserID(ctx context.Context, userID uint64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.WithContext(ctx).
		Table("menus").
		Select("DISTINCT menus.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("menus.deleted_at IS NULL").
		Where("menus.status = 1").
		Order("menus.parent_id ASC, menus.sort ASC, menus.id ASC").
		Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("list menus by user id: %w", err)
	}

	return menus, nil
}
