package model

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单主表
type Order struct {
	gorm.Model
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	OrderNum        string     `gorm:"unique;" json:"order_num"`          // 订单号，推荐用雪花算法
	AllPrice        int64      `json:"all_price"`                         // 订单总价，单位：分
	PayStatus       int        `gorm:"default:0" json:"pay_status"`       // 0:未支付 1:已支付
	PayType         int        `json:"pay_type"`                          // 1:支付宝 2:微信
	PayTime         *time.Time `json:"pay_time"`                          // 支付时间（指针类型，允许 NULL）
	TradeNo         string     `gorm:"type:varchar(64)" json:"trade_no"`  // 支付平台交易流水号（支付宝/微信返回）
	OrderStatus     int        `gorm:"default:0" json:"order_status"`     // 0:创建 1:支付 2:发货 3:完成 4:取消
	Type            int        `json:"type"`                              // 1:普通订单 2:秒杀订单
	AddressSnapshot string     `gorm:"type:text" json:"address_snapshot"` // 收货地址快照（JSON格式）
	ExpireTime      time.Time  `json:"expire_time"`                       // 订单过期时间（用于自动关单）
	Remark          string     `gorm:"type:text" json:"remark"`           // 用户备注
	TrackingNumber  string     `json:"tracking_number"`                   // 物流单号
	ShipTime        *time.Time `json:"ship_time"`                         // 发货时间（指针类型，允许 NULL）
	FinishTime      *time.Time `json:"finish_time"`                       // 完成时间（指针类型，允许 NULL）
	CancelTime      *time.Time `json:"cancel_time"`                       // 取消时间（指针类型，允许 NULL）
	AdminRemark     string     `gorm:"type:text" json:"admin_remark"`     // 管理员备注
	Version         int        `gorm:"default:0" json:"version"`          // 乐观锁版本号
}

// OrderItem 订单详情表 (商品快照)
type OrderItem struct {
	gorm.Model
	OrderNum     string `gorm:"index;not null" json:"order_num"` // 改为普通索引，一个订单可以有多个商品
	ProductID    uint   `json:"product_id"`
	ProductSkuID uint   `json:"product_sku_id"`
	Num          int    `json:"num"`      // 购买数量
	Price        int64  `json:"price"`    // 购买时的单价，单位：分（重点！）
	Title        string `json:"title"`    // 购买时的商品名
	ImgPath      string `json:"img_path"` // 购买时的图片
}
