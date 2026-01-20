package dao

import "xiaomi-mall/internal/model"

type addressDao struct{}

var Address = new(addressDao)

// 根据用户id查询地址列表
func (d *addressDao) GetbyUserId(userId uint) (address []*model.Address, err error) {
	err = DB.Where("user_id =?", userId).Order("is_default DESC").Find(&address).Error
	return
}

// 根据地址id查询地址
func (d *addressDao) GetbyAddressId(addressId uint) (address *model.Address, err error) {
	err = DB.Where("id =?", addressId).First(&address).Error
	return
}
