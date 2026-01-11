package model

import (
	"time"
	"gorm.io/gorm"
)

// SeckillActivity 秒杀活动场次表
// 例如：每日 10:00 场, 12:00 场
type SeckillActivity struct {
	gorm.Model
	Name      string    `gorm:"size:255;not null" json:"name"` // 活动名称 "双11早场"
	StartTime time.Time `gorm:"not null" json:"start_time"`    // 开始时间
	EndTime   time.Time `gorm:"not null" json:"end_time"`      // 结束时间
	Status    int       `gorm:"default:0" json:"status"`       // 0:未上线 1:已上线 2:已下线(人工干预)
	ImgPath   string    `json:"img_path"`                      // 活动宣传图
}

// SeckillProduct 秒杀商品表 (活动与商品的关联)
// 核心：高并发下，库存是热点数据，这里是最终一致性的落地表
type SeckillProduct struct {
	gorm.Model
	ActivityID   uint  `gorm:"not null;index" json:"activity_id"`     // 关联活动
	ProductID    uint  `gorm:"not null;index" json:"product_id"`      // 冗余SPU ID，方便列表展示
	ProductSkuID uint  `gorm:"not null;index" json:"product_sku_id"`  // 关联具体SKU
	
	SeckillPrice int64 `gorm:"not null" json:"seckill_price"`         // 秒杀价 (分)
	InitStock    int   `gorm:"not null" json:"init_stock"`            // 初始库存 (用于计算进度条)
	Stock        int   `gorm:"not null;check:stock>=0" json:"stock"`  // 剩余库存 (乐观锁扣减对象)
	
	LimitNum     int   `gorm:"default:1" json:"limit_num"`            // 每人限购 N 件
	Sort         int   `gorm:"default:0" json:"sort"`                 // 排序权重
	Version      int   `gorm:"default:0" json:"version"`              // 乐观锁版本号 (CAS机制)
}

// SeckillOrder 秒杀成功记录表
// 作用：1. 防止重复购买(唯一索引) 2. 异步削峰的落脚点 3. 对账
type SeckillOrder struct {
	gorm.Model
	// 联合唯一索引：确保 同一个用户 在 同一场活动 对 同一个商品 只能抢一次
	// 注意：uniqueIndex 的名字 idx_user_activity_product 必须一致，才能形成联合约束
	UserID           uint   `gorm:"not null;uniqueIndex:idx_user_activity_product" json:"user_id"`
	ActivityID       uint   `gorm:"not null;uniqueIndex:idx_user_activity_product" json:"activity_id"`
	SeckillProductID uint   `gorm:"not null;uniqueIndex:idx_user_activity_product" json:"seckill_product_id"`
	
	SkuID            uint   `json:"sku_id"`                        // 冗余记录 SKU ID
	Money            int64  `json:"money"`                         // 支付金额
	OrderNum         string `gorm:"type:varchar(32);index" json:"order_num"` // 关联主订单表的订单号
	Status           int    `gorm:"default:0" json:"status"`       // 0:待支付 1:已支付 2:超时取消
}