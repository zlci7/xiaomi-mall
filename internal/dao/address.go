package dao

import (
	"xiaomi-mall/internal/model"
)

var Address = new(AddressDao)

type AddressDao struct{}

// ========== 查询用户的所有地址 ==========
func (d *AddressDao) GetUserAddresses(userID uint) ([]*model.Address, error) {
	var addresses []*model.Address
	err := DB.Where("user_id = ?", userID).
		Order("is_default DESC, id DESC"). // 默认地址排在前面
		Find(&addresses).Error
	return addresses, err
}

// ========== 根据 ID 查询地址 ==========
func (d *AddressDao) GetByID(addressID uint) (*model.Address, error) {
	var address model.Address
	err := DB.Where("id = ?", addressID).First(&address).Error
	return &address, err
}

// ========== 创建地址 ==========
func (d *AddressDao) CreateAddress(address *model.Address) error {
	return DB.Create(address).Error
}

// ========== 更新地址 ==========
func (d *AddressDao) UpdateAddress(address *model.Address) error {
	return DB.Save(address).Error
}

// ========== 删除地址 ==========
func (d *AddressDao) DeleteAddress(addressID uint) error {
	return DB.Delete(&model.Address{}, addressID).Error
}

// ========== 取消所有默认地址（用于设置新的默认地址）==========
func (d *AddressDao) CancelDefaultAddress(userID uint) error {
	return DB.Model(&model.Address{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}

// ========== 查询用户的默认地址 ==========
func (d *AddressDao) GetDefaultAddress(userID uint) (*model.Address, error) {
	var address model.Address
	err := DB.Where("user_id = ? AND is_default = ?", userID, true).
		First(&address).Error
	return &address, err
}

// ========== 统计用户的地址数量 ==========
func (d *AddressDao) CountUserAddresses(userID uint) (int64, error) {
	var count int64
	err := DB.Model(&model.Address{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
