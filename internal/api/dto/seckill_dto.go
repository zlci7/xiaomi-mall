package dto

// ========== 秒杀活动列表 ==========
type SeckillActivityListReq struct {
	Page     int  `form:"page" binding:"omitempty,min=1"`
	PageSize int  `form:"page_size" binding:"omitempty,min=1,max=50"`
	Status   *int `form:"status" binding:"omitempty,oneof=0 1 2"` // 0:未上线 1:已上线 2:已下线
}

// ========== 秒杀活动详情 ==========
type SeckillActivityDetailReq struct {
	ActivityID uint `uri:"activity_id" binding:"required,min=1"`
}

// ========== 秒杀下单 ==========
type CreateSeckillOrderReq struct {
	ActivityID       uint `json:"activity_id" binding:"required,min=1"`
	SeckillProductID uint `json:"seckill_product_id" binding:"required,min=1"`
	Num              int  `json:"num" binding:"required,min=1"` // 购买数量
	AddressID        uint `json:"address_id" binding:"required,min=1"`
}
