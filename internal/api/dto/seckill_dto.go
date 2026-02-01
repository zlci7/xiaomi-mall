package dto

// ==================== 管理端：秒杀商品管理 ====================

// CreateSeckillProductReq 创建秒杀商品
type CreateSeckillProductReq struct {
	ProductID    uint   `json:"product_id" binding:"required,min=1"`
	SkuID        uint   `json:"sku_id" binding:"required,min=1"`
	SeckillPrice uint   `json:"seckill_price" binding:"required,min=1"` // 秒杀价（单位：分）
	SeckillStock uint   `json:"seckill_stock" binding:"required,min=1"` // 秒杀库存
	StartTime    string `json:"start_time" binding:"required"`          // 格式："2026-01-23 10:00:00"
	EndTime      string `json:"end_time" binding:"required"`            // 格式："2026-01-23 12:00:00"
}

// DeleteSeckillProductReq 删除秒杀商品
type DeleteSeckillProductReq struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateSeckillStatusReq 手动开始/结束秒杀
type UpdateSeckillStatusReq struct {
	ID     uint `uri:"id" binding:"required,min=1"`   // 路径参数
	Status int  `form:"status" binding:"oneof=0 1 2"` // 查询参数：0:未开始 1:进行中 2:已结束（去掉 required）
}

// SeckillProductListReq 秒杀商品列表（管理端）
type SeckillProductListReq struct {
	Page     int   `form:"page" binding:"omitempty,min=1"`
	PageSize int   `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   *int8 `form:"status" binding:"omitempty,oneof=0 1 2"` // 支持状态筛选
}

// SeckillProductDetailReq 秒杀商品详情
type SeckillProductDetailReq struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// PreheatSeckillProductReq 预热秒杀商品到 Redis
type PreheatSeckillProductReq struct {
	ID uint `uri:"id" binding:"required,min=1"` // 秒杀商品 ID
}

// ==================== 用户端：秒杀下单 ====================

// CreateSeckillOrderReq 用户秒杀下单
type CreateSeckillOrderReq struct {
	SeckillProductID uint `json:"seckill_product_id" binding:"required,min=1"`
	AddressID        uint `json:"address_id" binding:"required,min=1"` // 收货地址
}

// ==================== 用户端：秒杀查询 ====================

// UserSeckillListReq 用户端秒杀商品列表
type UserSeckillListReq struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// UserSeckillDetailReq 用户端秒杀商品详情
type UserSeckillDetailReq struct {
	ID uint `uri:"id" binding:"required,min=1"` // 秒杀商品 ID
}
