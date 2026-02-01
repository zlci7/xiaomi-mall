package vo

import "time"

// ==================== 管理端：秒杀商品管理 ====================

// CreateSeckillProductResp 创建秒杀商品响应
type CreateSeckillProductResp struct {
	ID uint `json:"id"`
}

// ==================== 用户端：秒杀列表 ====================

// UserSeckillListItemVO 秒杀商品列表项（简化信息）
type UserSeckillListItemVO struct {
	// 基础信息
	ID          uint   `json:"id"`           // 秒杀商品 ID
	ProductID   uint   `json:"product_id"`   // 商品 ID
	ProductName string `json:"product_name"` // 商品名称
	ImgPath     string `json:"img_path"`     // 商品图片

	// 价格信息
	OriginalPrice uint `json:"original_price"` // 原价（分）
	SeckillPrice  uint `json:"seckill_price"`  // 秒杀价（分）

	// 库存信息
	SeckillStock uint `json:"seckill_stock"` // 剩余库存
	SoldNum      uint `json:"sold_num"`      // 已售数量

	// 时间信息
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间

	// 状态信息
	Status    string `json:"status"`      // 状态："未开始" | "进行中" | "已结束"
	IsSoldOut bool   `json:"is_sold_out"` // 是否售罄
	CanBuy    bool   `json:"can_buy"`     // 是否可购买
}

// UserSeckillListResp 秒杀列表响应
type UserSeckillListResp struct {
	List  []UserSeckillListItemVO `json:"list"`
	Total int64                   `json:"total"`
}

// ==================== 用户端：秒杀详情 ====================

// UserSeckillDetailVO 秒杀商品详情（完整信息）
type UserSeckillDetailVO struct {
	// 基础信息
	ID           uint   `json:"id"`
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductTitle string `json:"product_title"` // 商品标题
	ProductInfo  string `json:"product_info"`  // 商品描述
	ImgPath      string `json:"img_path"`      // 主图
	CategoryID   uint   `json:"category_id"`   // 分类 ID

	// SKU 信息
	SkuID    uint   `json:"sku_id"`
	SkuTitle string `json:"sku_title"` // 规格名称
	SkuCode  string `json:"sku_code"`  // 商家编码

	// 价格信息
	OriginalPrice uint `json:"original_price"` // 原价（分）
	SeckillPrice  uint `json:"seckill_price"`  // 秒杀价（分）

	// 库存信息
	SeckillStock uint `json:"seckill_stock"` // 剩余库存
	TotalStock   uint `json:"total_stock"`   // 总库存
	SoldNum      uint `json:"sold_num"`      // 已售数量

	// 时间信息
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间

	// 状态信息
	Status       string `json:"status"`        // 状态："未开始" | "进行中" | "已结束"
	IsSoldOut    bool   `json:"is_sold_out"`   // 是否售罄
	CanBuy       bool   `json:"can_buy"`       // 是否可购买
	HasPurchased bool   `json:"has_purchased"` // 当前用户是否已购买
}

// ==================== 秒杀下单响应 ====================

// CreateSeckillOrderResp 秒杀下单响应
type CreateSeckillOrderResp struct {
	OrderNum    string    `json:"order_num"`
	TotalAmount int64     `json:"total_amount"`
	ExpireTime  time.Time `json:"expire_time"` // 订单过期时间（30分钟后）
	PayUrl      string    `json:"pay_url"`     // 支付链接（可选）
}
