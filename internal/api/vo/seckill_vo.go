package vo

import "time"

// ==================== 管理端：秒杀商品管理 ====================

// CreateSeckillProductResp 创建秒杀商品响应
type CreateSeckillProductResp struct {
	ID uint `json:"id"`
}

// SeckillProductItemVO 秒杀商品列表项（管理端）
type SeckillProductItemVO struct {
	ID            uint      `json:"id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"` // 商品名称
	SkuID         uint      `json:"sku_id"`
	SkuTitle      string    `json:"sku_title"`      // SKU规格
	ImgPath       string    `json:"img_path"`       // 商品图片
	OriginalPrice uint      `json:"original_price"` // 原价（单位：分）
	SeckillPrice  uint      `json:"seckill_price"`  // 秒杀价（单位：分）
	SeckillStock  uint      `json:"seckill_stock"`  // 剩余库存
	TotalStock    uint      `json:"total_stock"`    // 总库存
	SoldNum       uint      `json:"sold_num"`       // 已售数量
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        int8      `json:"status"`      // 0:未开始 1:进行中 2:已结束
	StatusText    string    `json:"status_text"` // 状态文本
	CreatedAt     time.Time `json:"created_at"`
}

// SeckillProductListResp 秒杀商品列表响应（管理端）
type SeckillProductListResp struct {
	List  []SeckillProductItemVO `json:"list"`
	Total int64                  `json:"total"`
}

// SeckillProductDetailResp 秒杀商品详情响应（管理端）
type SeckillProductDetailResp struct {
	ID            uint      `json:"id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductTitle  string    `json:"product_title"` // 商品副标题
	SkuID         uint      `json:"sku_id"`
	SkuTitle      string    `json:"sku_title"`
	SkuCode       string    `json:"sku_code"` // SKU编码
	ImgPath       string    `json:"img_path"`
	OriginalPrice uint      `json:"original_price"`
	SeckillPrice  uint      `json:"seckill_price"`
	SeckillStock  uint      `json:"seckill_stock"` // 剩余库存
	TotalStock    uint      `json:"total_stock"`   // 总库存（创建时的库存）
	SoldNum       uint      `json:"sold_num"`      // 已售数量
	Discount      string    `json:"discount"`      // 折扣（如："7.5折"）
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        int8      `json:"status"`
	StatusText    string    `json:"status_text"`
	TimeStatus    string    `json:"time_status"` // "即将开始" "进行中" "已结束"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ==================== 用户端：秒杀商品展示 ====================

// UserSeckillProductVO 用户端秒杀商品（列表和详情）
type UserSeckillProductVO struct {
	ID            uint      `json:"id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductTitle  string    `json:"product_title"`
	SkuID         uint      `json:"sku_id"`
	SkuTitle      string    `json:"sku_title"`
	ImgPath       string    `json:"img_path"`
	OriginalPrice uint      `json:"original_price"`
	SeckillPrice  uint      `json:"seckill_price"`
	SeckillStock  uint      `json:"seckill_stock"` // 剩余库存（可能从Redis读取）
	SoldNum       uint      `json:"sold_num"`      // 已售数量
	Discount      string    `json:"discount"`      // 折扣
	Progress      int       `json:"progress"`      // 进度百分比（0-100）
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsSoldOut     bool      `json:"is_sold_out"`   // 是否售罄
	CanBuy        bool      `json:"can_buy"`       // 当前用户是否可以购买（已购买则false）
	CountdownSec  int64     `json:"countdown_sec"` // 倒计时秒数（>0开始倒计时，<0已开始）
}

// UserSeckillListResp 用户端秒杀列表响应
type UserSeckillListResp struct {
	List  []UserSeckillProductVO `json:"list"`
	Total int64                  `json:"total"`
}

// ==================== 秒杀下单响应 ====================

// CreateSeckillOrderResp 秒杀下单响应
type CreateSeckillOrderResp struct {
	OrderNum    string    `json:"order_num"`
	TotalAmount int64     `json:"total_amount"`
	ExpireTime  time.Time `json:"expire_time"` // 订单过期时间（30分钟后）
	PayUrl      string    `json:"pay_url"`     // 支付链接（可选）
}
