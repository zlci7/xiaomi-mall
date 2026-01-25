package model

import (
	"time"

	"gorm.io/gorm"
)

// SeckillProduct 表：轻量级秒杀商品配置
type SeckillProduct struct {
	gorm.Model
	ProductID    uint      `gorm:"not null;index" json:"product_id"`
	SkuID        uint      `gorm:"not null;index" json:"sku_id"`
	SeckillPrice uint      `gorm:"not null" json:"seckill_price"`    // 秒杀价（单位：分）
	SeckillStock uint      `gorm:"not null" json:"seckill_stock"`    // 剩余库存
	TotalStock   uint      `gorm:"not null" json:"total_stock"`      // 总库存（初始库存，用于计算已售）
	StartTime    time.Time `gorm:"not null;index" json:"start_time"` //格式"2026-01-23 10:00:00"
	EndTime      time.Time `gorm:"not null;index" json:"end_time"`
	Status       int8      `gorm:"default:0;index" json:"status"` // 0:未开始 1:进行中 2:已结束
	Version      int       `gorm:"default:0" json:"version"`      // 乐观锁
}

// SeckillOrder 表：秒杀订单记录（防重复购买）
type SeckillOrder struct {
	gorm.Model
	UserID           uint   `gorm:"not null;uniqueIndex:idx_user_seckill" json:"user_id"`
	SeckillProductID uint   `gorm:"not null;uniqueIndex:idx_user_seckill" json:"seckill_product_id"`
	OrderNum         string `gorm:"not null;index" json:"order_num"`
	Status           int8   `gorm:"default:0" json:"status"` // 0:待支付 1:已支付 2:已取消
}
