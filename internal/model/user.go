package model

import "time"

type User struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	Username    string     `gorm:"column:username"`
	Password    string     `gorm:"column:password"`
	Nickname    string     `gorm:"column:nickname"`
	Email       string     `gorm:"column:email"`
	Phone       string     `gorm:"column:phone"`
	Avatar      string     `gorm:"column:avatar"`
	Status      int        `gorm:"column:status"`
	Remark      string     `gorm:"column:remark"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
