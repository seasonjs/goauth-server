package repository

import (
	"errors"
	"gorm.io/gorm"
	"oauthServer/entity"
	"strings"
)

//UserRepository 用户
type UserRepository interface {
	Save(*entity.User) (*entity.User, map[string]string)
	FindUserById(id int) (*entity.User, error)
	FindUsers() ([]entity.User, error)
	FindUserByName(string) (*entity.User, error)
	FindUserByNameAndPassword(username string, password string) (*entity.User, error)
	//GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

// UserRepo 通过*gorm.DB支持的mysql/等实现
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据库实例
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

//////UserRepo implements the repository.UserRepository interface
//var _ UserRepository = &UserRepo{}

func (r *UserRepo) Save(user *entity.User) (*entity.User, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *UserRepo) FindUserById(id int) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) FindUserByName(name string) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("username = ?", name).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) FindUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Debug().Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepo) FindUserByNameAndPassword(username string, password string) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where(&entity.User{Username: username, Password: password}).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
