// pkg/constants/order_status.go
package constants

const (
	// OrderStatus 订单状态
	ORDER_STATUS_PENDING   = 0 // 待支付
	ORDER_STATUS_PAID      = 1 // 已支付
	ORDER_STATUS_SHIPPED   = 2 // 已发货
	ORDER_STATUS_COMPLETED = 3 // 已完成
	ORDER_STATUS_CANCELLED = 4 // 已取消

	// PayStatus 支付状态
	PAY_STATUS_UNPAID = 0 // 未支付
	PAY_STATUS_PAID   = 1 // 已支付
	PAY_STATUS_REFUND = 2 // 已退款

	// PayType 支付方式
	PAY_TYPE_ALIPAY = 1 // 支付宝
	PAY_TYPE_WECHAT = 2 // 微信

	// OrderType 订单类型
	ORDER_TYPE_NORMAL  = 1 // 普通订单
	ORDER_TYPE_SECKILL = 2 // 秒杀订单
)
