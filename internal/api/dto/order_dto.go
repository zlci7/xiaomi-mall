package dto

// ========== 创建订单 ==========
type CreateOrderReq struct {
	Items     []OrderItemReq `json:"items" binding:"required,min=1,dive"` // 购买的商品列表
	AddressID uint           `json:"address_id" binding:"required,min=1"` // 收货地址 ID
	Remark    string         `json:"remark" binding:"omitempty,max=200"`  // 用户备注
}

// 订单商品项
type OrderItemReq struct {
	SkuID uint `json:"sku_id" binding:"required,min=1"` // SKU ID
	Num   int  `json:"num" binding:"required,min=1"`    // 购买数量
}

// ========== 支付订单 ==========
type PayOrderReq struct {
	OrderNo string `json:"order_no" binding:"required"`           // 订单号
	PayType int    `json:"pay_type" binding:"required,oneof=1 2"` // 1:支付宝 2:微信
}

// ========== 取消订单 ==========
type CancelOrderReq struct {
	OrderNo string `json:"order_no" binding:"required"` // 订单号
}

// ========== 订单列表查询 ==========
type OrderListReq struct {
	Page        int  `form:"page" binding:"omitempty,min=1"`                   // 页码
	PageSize    int  `form:"page_size" binding:"omitempty,min=1,max=50"`       // 每页数量
	OrderStatus *int `form:"order_status" binding:"omitempty,oneof=0 1 2 3 4"` // 订单状态筛选（nil=全部）
}

// ========== 订单详情查询 ==========
type OrderDetailReq struct {
	OrderNo string `uri:"order_no" binding:"required"` // 路径参数
}

// ========== 确认收货 ==========
type ConfirmOrderReq struct {
	OrderNo string `json:"order_no" binding:"required"`
}

// ========== 管理端：发货 ==========
type ShipOrderReq struct {
	OrderNo        string `json:"order_no" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required,max=50"` // 物流单号
	AdminRemark    string `json:"admin_remark" binding:"omitempty,max=200"`  // 管理员备注
}
