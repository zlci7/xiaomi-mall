package dao

import (
	"time"
	"xiaomi-mall/internal/model"

	"gorm.io/gorm"
)

var Order = new(OrderDao)

type OrderDao struct{}

// ========== 创建订单（事务版本）==========
func (d *OrderDao) CreateOrder(tx *gorm.DB, order *model.Order, orderItems []*model.OrderItem) error {
	// 1️⃣ 创建订单主表
	if err := tx.Create(order).Error; err != nil {
		return err
	}

	// 2️⃣ 批量创建订单详情（性能优化）
	if len(orderItems) > 0 {
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}
	}

	return nil
}

// ========== 查询订单（根据订单号）==========
func (d *OrderDao) GetOrderByOrderNum(orderNum string) (*model.Order, error) {
	var order model.Order
	err := DB.Where("order_num = ?", orderNum).First(&order).Error
	return &order, err
}

// ========== 查询订单详情 ==========
func (d *OrderDao) GetOrderItems(orderNum string) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	err := DB.Where("order_num = ?", orderNum).Find(&items).Error
	return items, err
}

// ========== 更新订单状态（乐观锁）==========
func (d *OrderDao) UpdateOrderStatus(tx *gorm.DB, orderNum string, status int, version int) (int64, error) {
	result := tx.Model(&model.Order{}).
		Where("order_num = ? AND version = ?", orderNum, version).
		Updates(map[string]interface{}{
			"order_status": status,
			"version":      version + 1,
		})

	return result.RowsAffected, result.Error
}

// ========== 支付订单（乐观锁）==========
func (d *OrderDao) PayOrder(orderNum string, payType int, tradeNo string, version int) (int64, error) {
	result := DB.Model(&model.Order{}).
		Where("order_num = ? AND pay_status = 0 AND version = ?", orderNum, version).
		Updates(map[string]interface{}{
			"pay_status":   1,
			"order_status": 1,
			"pay_type":     payType,
			"pay_time":     time.Now(),
			"trade_no":     tradeNo,
			"version":      version + 1,
		})

	return result.RowsAffected, result.Error
}
