package vo

import "time"

// ========================================
// 管理端 - 秒杀活动管理响应
// ========================================

// ========== 创建/编辑秒杀活动响应 ==========
type CreateSeckillActivityResp struct {
	ActivityID uint      `json:"activity_id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     int       `json:"status"`
	ImgPath    string    `json:"img_path"`
}

// ========== 添加/编辑秒杀商品响应 ==========
type AddSeckillProductResp struct {
	SeckillProductID uint   `json:"seckill_product_id"`
	ActivityID       uint   `json:"activity_id"`
	ProductID        uint   `json:"product_id"`
	ProductSkuID     uint   `json:"product_sku_id"`
	ProductName      string `json:"product_name"`
	SkuTitle         string `json:"sku_title"`
	OriginalPrice    int64  `json:"original_price"`
	SeckillPrice     int64  `json:"seckill_price"`
	Stock            int    `json:"stock"`
	LimitNum         int    `json:"limit_num"`
	Sort             int    `json:"sort"`
}

// ========================================
// 用户端 - 秒杀活动查询响应
// ========================================

// ========== 秒杀活动列表响应 ==========
type SeckillActivityListResp struct {
	List     []SeckillActivityVO `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// 秒杀活动 VO
type SeckillActivityVO struct {
	ActivityID   uint      `json:"activity_id"`
	Name         string    `json:"name"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       int       `json:"status"`      // 0:未上线 1:已上线 2:已下线
	StatusText   string    `json:"status_text"` // 状态文本（前端展示）
	ImgPath      string    `json:"img_path"`
	ProductCount int       `json:"product_count"` // 该活动的商品数量

	// 时间状态（前端展示用）
	TimeStatus string `json:"time_status"` // "未开始" / "进行中" / "已结束"
}

// ========== 秒杀活动详情响应 ==========
type SeckillActivityDetailResp struct {
	ActivityID uint      `json:"activity_id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     int       `json:"status"`
	StatusText string    `json:"status_text"`
	ImgPath    string    `json:"img_path"`
	TimeStatus string    `json:"time_status"`

	// 该活动下的秒杀商品列表
	Products []SeckillProductVO `json:"products"`
}

// 秒杀商品 VO
type SeckillProductVO struct {
	SeckillProductID uint   `json:"seckill_product_id"`
	ProductID        uint   `json:"product_id"`
	ProductSkuID     uint   `json:"product_sku_id"`
	ProductName      string `json:"product_name"` // 从 Product 表查询
	ProductImg       string `json:"product_img"`  // 从 Product 表查询
	SkuTitle         string `json:"sku_title"`    // 从 ProductSku 表查询

	OriginalPrice int64 `json:"original_price"` // 原价（从 SKU 表）
	SeckillPrice  int64 `json:"seckill_price"`  // 秒杀价
	Discount      int   `json:"discount"`       // 折扣（如：25 表示2.5折）

	InitStock    int `json:"init_stock"`    // 初始库存
	Stock        int `json:"stock"`         // 剩余库存
	StockPercent int `json:"stock_percent"` // 剩余百分比（前端进度条）
	SoldNum      int `json:"sold_num"`      // 已售数量

	LimitNum  int  `json:"limit_num"` // 限购数量
	Sort      int  `json:"sort"`
	IsSoldOut bool `json:"is_sold_out"` // 是否售罄
}

// ========================================
// 用户端 - 秒杀下单响应
// ========================================

// ========== 秒杀下单响应 ==========
type CreateSeckillOrderResp struct {
	OrderNo     string    `json:"order_no"`
	TotalAmount int64     `json:"total_amount"`
	ExpireTime  time.Time `json:"expire_time"`
	PayUrl      string    `json:"pay_url"`
}
