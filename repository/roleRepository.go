package repository

import (
	"errors"
	"gorm.io/gorm"
	"oauthServer/entity"
)

type RoleRepository interface {
	Save(role *entity.Role) (*entity.Role, map[string]string)
	FindGroupByRoleType() ([]entity.Role, error)
	FindRoleByUserId(userId int) ([]entity.Role, error)
}

// RoleRepo 通过*gorm.DB支持的mysql/等实现
type RoleRepo struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色数据库实例
func NewRoleRepository(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db}
}

//////UserRepo implements the repository.UserRepository interface
//var _ UserRepository = &UserRepo{}

func (r *RoleRepo) Save(role *entity.Role) (*entity.Role, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&role).Error
	if err != nil {
		//If the email is already taken
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return role, nil
}

func (r *RoleRepo) FindGroupByRoleType() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Debug().Distinct("role_type,role_name").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return roles, nil
}

func (r *RoleRepo) FindRoleByUserId(userId int) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Debug().Where("user_id=?", userId).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return roles, nil
}
