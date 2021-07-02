package entity

import (
	"gorm.io/datatypes"
)

// Scope .
type Scope struct {
	changes    map[string]interface{}
	ID         int            `gorm:"primaryKey;column:id"`
	UserId     int            `gorm:"column:user_id"`
	ScopeName  string         `gorm:"column:scope_name"`
	ScopeDesc  string         `gorm:"column:scope_desc"`
	CreateTime datatypes.Date `gorm:"column:create_time"`
	UpdateTime datatypes.Date `gorm:"column:update_time"`
}

// TableName .
func (obj *Scope) TableName() string {
	return "scope"
}
