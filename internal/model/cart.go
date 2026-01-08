package model

import "gorm.io/gorm"

// Cart 购物车
// 说明：MySQL 用于持久化存储，Redis 用于读写缓存提升性能
type Cart struct {
	gorm.Model
	UserID       uint `gorm:"not null;uniqueIndex:idx_user_sku" json:"user_id"` // 联合唯一索引
	ProductID    uint `gorm:"not null" json:"product_id"`
	ProductSkuID uint `gorm:"not null;uniqueIndex:idx_user_sku" json:"product_sku_id"` // 联合唯一索引，防止重复添加
	Num          uint `gorm:"not null" json:"num"`                                     // 数量
	Check        bool `json:"check"`                                                   // 是否勾选
}
