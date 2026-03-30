package auth

import (
	"context"
	"errors"
	"fmt"

	"sys-admin-serve/internal/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrDuplicateUsername = errors.New("duplicate username")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(txRepo *Repository) error) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	}); err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	return nil
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

func (r *Repository) GetRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query role by code: %w", err)
	}

	return &role, nil
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

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateEntryError(err) {
			return ErrDuplicateUsername
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) CreateUserRole(ctx context.Context, userID, roleID uint64) error {
	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := r.db.WithContext(ctx).Create(&userRole).Error; err != nil {
		return fmt.Errorf("create user role: %w", err)
	}

	return nil
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
