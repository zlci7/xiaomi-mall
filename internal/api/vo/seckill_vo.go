package vo

import "time"

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
	Status       int       `json:"status"`
	ImgPath      string    `json:"img_path"`
	ProductCount int       `json:"product_count"` // 该活动的商品数量
}

// ========== 秒杀活动详情响应 ==========
type SeckillActivityDetailResp struct {
	ActivityID uint      `json:"activity_id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     int       `json:"status"`
	ImgPath    string    `json:"img_path"`

	// 该活动下的秒杀商品列表
	Products []SeckillProductVO `json:"products"`
}

// 秒杀商品 VO
type SeckillProductVO struct {
	SeckillProductID uint   `json:"seckill_product_id"`
	ProductID        uint   `json:"product_id"`
	ProductName      string `json:"product_name"` // 从 Product 表查询
	ProductImg       string `json:"product_img"`  // 从 Product 表查询
	SkuTitle         string `json:"sku_title"`    // 从 ProductSku 表查询

	OriginalPrice int64 `json:"original_price"` // 原价（从 SKU 表）
	SeckillPrice  int64 `json:"seckill_price"`  // 秒杀价

	InitStock    int `json:"init_stock"`    // 初始库存
	Stock        int `json:"stock"`         // 剩余库存
	StockPercent int `json:"stock_percent"` // 剩余百分比（前端进度条）

	LimitNum int `json:"limit_num"` // 限购数量
	Sort     int `json:"sort"`
}

// ========== 秒杀下单响应 ==========
type CreateSeckillOrderResp struct {
	OrderNo     string    `json:"order_no"`
	TotalAmount int64     `json:"total_amount"`
	ExpireTime  time.Time `json:"expire_time"`
	PayUrl      string    `json:"pay_url"`
}
