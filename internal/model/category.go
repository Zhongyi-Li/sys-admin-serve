package model

import "time"

type Category struct {
	ID        uint64     `gorm:"column:id;primaryKey"`
	ParentID  uint64     `gorm:"column:parent_id"`
	Name      string     `gorm:"column:name"`
	Code      string     `gorm:"column:code"`
	Sort      int        `gorm:"column:sort"`
	Status    int        `gorm:"column:status"`
	Icon      string     `gorm:"column:icon"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (Category) TableName() string {
	return "categories"
}
