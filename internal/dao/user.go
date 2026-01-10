package dao

import (
	"xiaomi-mall/internal/model"
)

var User = new(UserDao)

type UserDao struct{}

// ExistOrNotByUserName 判断用户名是否存在
func (d *UserDao) ExistOrNotByUserName(userName string) (exist bool, err error) {
	var count int64
	err = DB.Model(&model.User{}).Where("user_name = ?", userName).Count(&count).Error
	if count > 0 {
		return true, err
	}
	return false, err
}

// CreateUser 创建用户
func (d *UserDao) CreateUser(user *model.User) error {
	return DB.Model(&model.User{}).Create(user).Error
}

// GetUserByUserName 根据用户名获取用户
func (d *UserDao) GetUserByUserName(userName string) (user *model.User, err error) {
	err = DB.Model(&model.User{}).Where("user_name = ?", userName).First(&user).Error
	return
}
