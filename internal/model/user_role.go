package model

import "time"

type UserRole struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	UserID    uint64    `gorm:"column:user_id"`
	RoleID    uint64    `gorm:"column:role_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
