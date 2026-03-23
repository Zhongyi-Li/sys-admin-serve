package model

import "time"

type Role struct {
	ID        uint64     `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name"`
	Code      string     `gorm:"column:code"`
	Status    int        `gorm:"column:status"`
	Sort      int        `gorm:"column:sort"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}
