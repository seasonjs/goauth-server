package entity

import (
	"gorm.io/datatypes"
)

// Client .
type Client struct {
	changes      map[string]interface{}
	ID           int            `gorm:"primaryKey;column:id"`
	UserId       int            `gorm:"column:user_id"`
	ClientName   string         `gorm:"column:client_name"`
	ClientKey    string         `gorm:"column:client_key"`
	ClientSecret string         `gorm:"column:client_secret"`
	CreateTime   datatypes.Date `gorm:"column:create_time"`
	UpdateTime   datatypes.Date `gorm:"column:update_time"`
}

// TableName .
func (obj *Client) TableName() string {
	return "client"
}
