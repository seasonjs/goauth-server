package repository

import (
	"errors"
	"gorm.io/gorm"
	"oauthServer/entity"
)

type PermissionRepository interface {
	Save(permission *entity.Permission) (*entity.Permission, map[string]string)
	FindPermissionByRoleType(rt int) ([]entity.Permission, error)
	FindPermissionByPermissionResourceAndRoleType(resource string, roleType int) (*entity.Permission, error)
}

// PermissionRepo 通过*gorm.DB支持的mysql/等实现
type PermissionRepo struct {
	db *gorm.DB
}

// NewPermissionRepository 创建角色数据库实例
func NewPermissionRepository(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db}
}

func (r *PermissionRepo) Save(permission *entity.Permission) (*entity.Permission, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&permission).Error
	if err != nil {
		//If the email is already taken
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return permission, nil
}
func (r *PermissionRepo) FindPermissionByRoleType(rt int) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Debug().Where("role_type=?", rt).Find(&permissions).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("permission not found")
	}
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
func (r *PermissionRepo) FindPermissionByPermissionResourceAndRoleType(resource string, roleType int) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.Debug().Where(&entity.Permission{PermissionResource: resource, RoleType: roleType}).Take(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("permission not found")
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
