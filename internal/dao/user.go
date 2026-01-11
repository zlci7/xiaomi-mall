package dao

import (
	"xiaomi-mall/internal/model"
)

var User = new(UserDao)

type UserDao struct{}

// ExistOrNotByPhone 判断手机号是否存在
func (d *UserDao) ExistOrNotByPhone(phone string) (exist bool, err error) {
	var count int64
	err = DB.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	if count > 0 {
		return true, err
	}
	return false, err
}

// CreateUser 创建用户
func (d *UserDao) CreateUser(user *model.User) error {
	return DB.Model(&model.User{}).Create(user).Error
}

// GetUserByPhone 根据手机号获取用户
func (d *UserDao) GetUserByPhone(phone string) (user *model.User, err error) {
	err = DB.Model(&model.User{}).Where("phone = ?", phone).First(&user).Error
	return
}
