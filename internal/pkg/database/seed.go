package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const seedTimeout = 10 * time.Second

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "Admin@123456"
	defaultAdminNickname = "System Administrator"
	defaultAdminRoleName = "Super Admin"
	defaultAdminRoleCode = "super_admin"

	systemMenuName     = "system"
	userManagementName = "user_management"
	roleManagementName = "role_management"
)

type userSeedRecord struct {
	ID        uint64    `gorm:"column:id"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Nickname  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	Avatar    string    `gorm:"column:avatar"`
	Status    int       `gorm:"column:status"`
	Remark    string    `gorm:"column:remark"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (userSeedRecord) TableName() string {
	return "users"
}

type roleSeedRecord struct {
	ID        uint64    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	Status    int       `gorm:"column:status"`
	Sort      int       `gorm:"column:sort"`
	Remark    string    `gorm:"column:remark"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (roleSeedRecord) TableName() string {
	return "roles"
}

type userRoleSeedRecord struct {
	ID        uint64    `gorm:"column:id"`
	UserID    uint64    `gorm:"column:user_id"`
	RoleID    uint64    `gorm:"column:role_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (userRoleSeedRecord) TableName() string {
	return "user_roles"
}

type menuSeedRecord struct {
	ID         uint64    `gorm:"column:id"`
	ParentID   uint64    `gorm:"column:parent_id"`
	Name       string    `gorm:"column:name"`
	Title      string    `gorm:"column:title"`
	Path       string    `gorm:"column:path"`
	Component  string    `gorm:"column:component"`
	Icon       string    `gorm:"column:icon"`
	MenuType   string    `gorm:"column:menu_type"`
	Permission string    `gorm:"column:permission"`
	Sort       int       `gorm:"column:sort"`
	Hidden     int       `gorm:"column:hidden"`
	Status     int       `gorm:"column:status"`
	Remark     string    `gorm:"column:remark"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (menuSeedRecord) TableName() string {
	return "menus"
}

type roleMenuSeedRecord struct {
	ID        uint64    `gorm:"column:id"`
	RoleID    uint64    `gorm:"column:role_id"`
	MenuID    uint64    `gorm:"column:menu_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (roleMenuSeedRecord) TableName() string {
	return "role_menus"
}

func SeedInitialData(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), seedTimeout)
	defer cancel()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roleID, err := ensureAdminRole(tx)
		if err != nil {
			return err
		}

		userID, err := ensureAdminUser(tx)
		if err != nil {
			return err
		}

		if err := ensureAdminUserRole(tx, userID, roleID); err != nil {
			return err
		}

		menuIDs, err := ensureSystemMenus(tx)
		if err != nil {
			return err
		}

		if err := ensureRoleMenus(tx, roleID, menuIDs); err != nil {
			return err
		}

		return nil
	})
}

func ensureAdminRole(tx *gorm.DB) (uint64, error) {
	now := time.Now()
	record := roleSeedRecord{
		Name:      defaultAdminRoleName,
		Code:      defaultAdminRoleCode,
		Status:    1,
		Sort:      1,
		Remark:    "seeded system role",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.Assignments(map[string]any{"name": record.Name, "status": record.Status, "sort": record.Sort, "remark": record.Remark, "updated_at": now}),
	}).Create(&record).Error; err != nil {
		return 0, fmt.Errorf("upsert admin role: %w", err)
	}

	var role roleSeedRecord
	if err := tx.Where("code = ?", defaultAdminRoleCode).First(&role).Error; err != nil {
		return 0, fmt.Errorf("query admin role: %w", err)
	}

	return role.ID, nil
}

func ensureAdminUser(tx *gorm.DB) (uint64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(resolveAdminPassword()), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash admin password: %w", err)
	}

	now := time.Now()
	record := userSeedRecord{
		Username:  resolveAdminUsername(),
		Password:  string(passwordHash),
		Nickname:  defaultAdminNickname,
		Status:    1,
		Remark:    "seeded system administrator",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.Assignments(map[string]any{"nickname": record.Nickname, "status": record.Status, "remark": record.Remark, "updated_at": now}),
	}).Create(&record).Error; err != nil {
		return 0, fmt.Errorf("upsert admin user: %w", err)
	}

	var user userSeedRecord
	if err := tx.Where("username = ?", record.Username).First(&user).Error; err != nil {
		return 0, fmt.Errorf("query admin user: %w", err)
	}

	return user.ID, nil
}

func ensureAdminUserRole(tx *gorm.DB, userID, roleID uint64) error {
	record := userRoleSeedRecord{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
		return fmt.Errorf("create admin user role: %w", err)
	}

	return nil
}

func ensureSystemMenus(tx *gorm.DB) ([]uint64, error) {
	now := time.Now()
	systemMenu := menuSeedRecord{
		ParentID:   0,
		Name:       systemMenuName,
		Title:      "System Management",
		Path:       "/system",
		Component:  "Layout",
		Icon:       "setting",
		MenuType:   "catalog",
		Permission: "",
		Sort:       10,
		Hidden:     0,
		Status:     1,
		Remark:     "seeded root menu",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := upsertMenuByName(tx, systemMenu, now); err != nil {
		return nil, err
	}

	storedSystemMenu, err := queryMenuByName(tx, systemMenuName)
	if err != nil {
		return nil, err
	}

	userManagementMenu := menuSeedRecord{
		ParentID:   storedSystemMenu.ID,
		Name:       userManagementName,
		Title:      "User Management",
		Path:       "/system/users",
		Component:  "system/users/index",
		Icon:       "user",
		MenuType:   "menu",
		Permission: "system:user:list",
		Sort:       11,
		Hidden:     0,
		Status:     1,
		Remark:     "seeded child menu",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := upsertMenuByName(tx, userManagementMenu, now); err != nil {
		return nil, err
	}

	roleManagementMenu := menuSeedRecord{
		ParentID:   storedSystemMenu.ID,
		Name:       roleManagementName,
		Title:      "Role Management",
		Path:       "/system/roles",
		Component:  "system/roles/index",
		Icon:       "safety-certificate",
		MenuType:   "menu",
		Permission: "system:role:list",
		Sort:       12,
		Hidden:     0,
		Status:     1,
		Remark:     "seeded child menu",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := upsertMenuByName(tx, roleManagementMenu, now); err != nil {
		return nil, err
	}

	storedUserMenu, err := queryMenuByName(tx, userManagementName)
	if err != nil {
		return nil, err
	}

	storedRoleMenu, err := queryMenuByName(tx, roleManagementName)
	if err != nil {
		return nil, err
	}

	return []uint64{storedSystemMenu.ID, storedUserMenu.ID, storedRoleMenu.ID}, nil
}

func upsertMenuByName(tx *gorm.DB, menu menuSeedRecord, now time.Time) error {
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"parent_id":  menu.ParentID,
			"title":      menu.Title,
			"path":       menu.Path,
			"component":  menu.Component,
			"icon":       menu.Icon,
			"menu_type":  menu.MenuType,
			"permission": menu.Permission,
			"sort":       menu.Sort,
			"hidden":     menu.Hidden,
			"status":     menu.Status,
			"remark":     menu.Remark,
			"updated_at": now,
		}),
	}).Create(&menu).Error; err != nil {
		return fmt.Errorf("upsert menu %s: %w", menu.Name, err)
	}

	return nil
}

func queryMenuByName(tx *gorm.DB, name string) (*menuSeedRecord, error) {
	var menu menuSeedRecord
	if err := tx.Where("name = ?", name).First(&menu).Error; err != nil {
		return nil, fmt.Errorf("query menu by name %s: %w", name, err)
	}

	return &menu, nil
}

func ensureRoleMenus(tx *gorm.DB, roleID uint64, menuIDs []uint64) error {
	now := time.Now()
	for _, menuID := range menuIDs {
		record := roleMenuSeedRecord{
			RoleID:    roleID,
			MenuID:    menuID,
			CreatedAt: now,
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
			return fmt.Errorf("create role menu relation role=%d menu=%d: %w", roleID, menuID, err)
		}
	}

	return nil
}

func resolveAdminUsername() string {
	value := os.Getenv("SEED_ADMIN_USERNAME")
	if value == "" {
		return defaultAdminUsername
	}

	return value
}

func resolveAdminPassword() string {
	value := os.Getenv("SEED_ADMIN_PASSWORD")
	if value == "" {
		return defaultAdminPassword
	}

	return value
}
