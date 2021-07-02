package entity

import (
	"gorm.io/datatypes"
)

// User .
type User struct {
	changes     map[string]interface{}
	ID          int            `gorm:"primaryKey;column:id"`
	CreateTime  datatypes.Date `gorm:"column:create_time"`
	UpdateTime  datatypes.Date `gorm:"column:update_time"`
	Username    string         `gorm:"column:username"`
	Password    string         `gorm:"column:password"`
	ExtendsInfo datatypes.JSON `gorm:"column:extends_info"`
	ScopeID     int            `gorm:"column:scope_id"`
	ClientID    int            `gorm:"column:client_id"`
}

// TableName .
func (obj *User) TableName() string {
	return "user"
}
