package vo

import "time"

// ========== 创建订单响应 ==========
type CreateOrderResp struct {
	OrderNo     string    `json:"order_no"`     // 订单号
	TotalAmount int64     `json:"total_amount"` // 订单总金额（分）
	ExpireTime  time.Time `json:"expire_time"`  // 过期时间
	PayUrl      string    `json:"pay_url"`      // 支付链接（可选）
}

// ========== 订单列表响应 ==========
type OrderListResp struct {
	List     []OrderItemVO `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// 订单列表项（简化版）
type OrderItemVO struct {
	OrderNo      string            `json:"order_no"`      // 订单号
	TotalAmount  int64             `json:"total_amount"`  // 订单总金额（分）
	OrderStatus  int               `json:"order_status"`  // 订单状态
	PayStatus    int               `json:"pay_status"`    // 支付状态
	ProductCount int               `json:"product_count"` // 商品总数
	FirstProduct ProductSnapshotVO `json:"first_product"` // 第一个商品（用于列表展示）
	CreatedAt    time.Time         `json:"created_at"`    // 创建时间
	ExpireTime   time.Time         `json:"expire_time"`   // 过期时间
}

// 商品快照（列表用）
type ProductSnapshotVO struct {
	Title   string `json:"title"`    // 商品名
	ImgPath string `json:"img_path"` // 图片
	Price   int64  `json:"price"`    // 单价（分）
	Num     int    `json:"num"`      // 数量
}

// ========== 订单详情响应 ==========
type OrderDetailResp struct {
	// 订单基本信息
	OrderNo     string `json:"order_no"`
	TotalAmount int64  `json:"total_amount"`
	OrderStatus int    `json:"order_status"`
	PayStatus   int    `json:"pay_status"`
	PayType     int    `json:"pay_type"`
	Type        int    `json:"type"` // 1:普通订单 2:秒杀订单

	// 时间信息
	CreatedAt  time.Time  `json:"created_at"`
	PayTime    *time.Time `json:"pay_time,omitempty"`
	ShipTime   *time.Time `json:"ship_time,omitempty"`
	FinishTime *time.Time `json:"finish_time,omitempty"`
	CancelTime *time.Time `json:"cancel_time,omitempty"`
	ExpireTime time.Time  `json:"expire_time"`

	// 收货地址（快照）
	Address AddressSnapshotVO `json:"address"`

	// 商品列表（快照）
	Items []OrderDetailItemVO `json:"items"`

	// 物流信息
	TrackingNumber string `json:"tracking_number,omitempty"` // 物流单号

	// 备注
	Remark      string `json:"remark,omitempty"`       // 用户备注
	AdminRemark string `json:"admin_remark,omitempty"` // 管理员备注（仅管理端）
}

// 收货地址快照
type AddressSnapshotVO struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

// 订单商品详情项
type OrderDetailItemVO struct {
	ProductID    uint   `json:"product_id"`
	ProductSkuID uint   `json:"product_sku_id"`
	Title        string `json:"title"`    // 商品名（快照）
	ImgPath      string `json:"img_path"` // 图片（快照）
	Price        int64  `json:"price"`    // 单价（快照，分）
	Num          int    `json:"num"`      // 数量
	Subtotal     int64  `json:"subtotal"` // 小计 = Price * Num
}

// ========== 支付订单响应 ==========
type PayOrderResp struct {
	OrderNo   string `json:"order_no"`
	PayStatus int    `json:"pay_status"` // 1:已支付
	TradeNo   string `json:"trade_no"`   // 支付平台交易流水号（模拟）
}
