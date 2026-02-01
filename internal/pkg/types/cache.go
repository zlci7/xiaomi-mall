package types

// SeckillProductCache Redis 缓存的秒杀商品数据结构
// 用于序列化/反序列化 Redis 中存储的秒杀商品详情
type SeckillProductCache struct {
	SeckillID     uint   `json:"seckill_id"`
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductTitle  string `json:"title"`
	ProductInfo   string `json:"info"`
	ProductImg    string `json:"img"`
	CategoryID    uint   `json:"category_id"`
	SkuID         uint   `json:"sku_id"`
	SkuTitle      string `json:"sku_title"`
	SkuCode       string `json:"sku_code"`
	OriginalPrice int64  `json:"original_price"`
	SeckillPrice  uint   `json:"seckill_price"`
	SeckillStock  uint   `json:"seckill_stock"` // 总库存（用于计算已售）
	TotalStock    uint   `json:"total_stock"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
}
