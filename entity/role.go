package entity

import (
	"gorm.io/datatypes"
)

// Role .
type Role struct {
	changes    map[string]interface{}
	ID         int            `gorm:"primaryKey;column:id"`
	UserId     int            `gorm:"column:user_id"`
	RoleType   int            `gorm:"column:role_type"`
	RoleName   string         `gorm:"column:role_name"`
	CreateTime datatypes.Date `gorm:"column:create_time"`
	UpdateTime datatypes.Date `gorm:"column:update_time"`
}

// TableName .
func (obj *Role) TableName() string {
	return "role"
}
