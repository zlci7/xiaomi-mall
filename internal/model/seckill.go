package model

import (
	"time"

	"gorm.io/gorm"
)

// SeckillActivity 秒杀活动表
type SeckillActivity struct {
	gorm.Model
	Name      string    `gorm:"size:255;not null" json:"name"` // 活动名称
	StartTime time.Time `gorm:"not null" json:"start_time"`    // 开始时间
	EndTime   time.Time `gorm:"not null" json:"end_time"`      // 结束时间
	Status    int       `gorm:"default:0" json:"status"`       // 0:未开始 1:进行中 2:已结束
}

// SeckillProduct 秒杀商品关联表
type SeckillProduct struct {
	gorm.Model
	ActivityID   uint  `gorm:"not null;index" json:"activity_id"`                    // 活动ID
	ProductSkuID uint  `gorm:"not null;index" json:"product_sku_id"`                 // 商品SKU ID
	SeckillPrice int64 `gorm:"not null" json:"seckill_price"`                        // 秒杀价格，单位：分
	SeckillStock int   `gorm:"not null;check:seckill_stock>=0" json:"seckill_stock"` // 秒杀库存
	LimitNum     int   `gorm:"default:1" json:"limit_num"`                           // 每人限购数量
	Version      int   `gorm:"default:0" json:"version"`                             // 乐观锁版本号
}

// SeckillOrder 秒杀订单记录表（用于防止用户重复购买）
type SeckillOrder struct {
	gorm.Model
	UserID           uint   `gorm:"not null;index:idx_user_activity_sku" json:"user_id"`
	ActivityID       uint   `gorm:"not null;index:idx_user_activity_sku" json:"activity_id"`
	SeckillProductID uint   `gorm:"not null;index:idx_user_activity_sku" json:"seckill_product_id"`
	OrderNum         string `gorm:"type:varchar(32)" json:"order_num"` // 关联的订单号
	Status           int    `gorm:"default:0" json:"status"`           // 0:待支付 1:已支付 2:已取消
}
