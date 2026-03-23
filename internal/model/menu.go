package model

import "time"

type Menu struct {
	ID         uint64     `gorm:"column:id;primaryKey"`
	ParentID   uint64     `gorm:"column:parent_id"`
	Name       string     `gorm:"column:name"`
	Title      string     `gorm:"column:title"`
	Path       string     `gorm:"column:path"`
	Component  string     `gorm:"column:component"`
	Icon       string     `gorm:"column:icon"`
	MenuType   string     `gorm:"column:menu_type"`
	Permission string     `gorm:"column:permission"`
	Sort       int        `gorm:"column:sort"`
	Hidden     int        `gorm:"column:hidden"`
	Status     int        `gorm:"column:status"`
	Remark     string     `gorm:"column:remark"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
}

func (Menu) TableName() string {
	return "menus"
}
