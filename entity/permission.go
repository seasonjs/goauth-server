package entity

// Permission .
type Permission struct {
	changes            map[string]interface{}
	ID                 int    `gorm:"primaryKey;column:id"`
	PermissionName     string `gorm:"column:permission_name"`
	PermissionResource string `gorm:"column:permission_resource"`
	PermissionDesc     string `gorm:"column:permission_desc"`
	RoleType           int    `gorm:"column:role_type"`
}

// TableName .
func (obj *Permission) TableName() string {
	return "permission"
}
