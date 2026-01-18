package router

import (
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// router/order_router.go
func RegisterOrderRoutes(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/order")
	orderGroup.Use(middleware.JWTAuth()) // JWT 认证
	{
		// 普通订单
		orderGroup.POST("/create", handler.CreateOrder)      // 创建订单
		orderGroup.POST("/pay", handler.PayOrder)            // 支付订单
		orderGroup.POST("/cancel", handler.CancelOrder)      // 取消订单
		orderGroup.GET("/list", handler.GetOrderList)        // 订单列表
		orderGroup.GET("/:order_no", handler.GetOrderDetail) // 订单详情
		orderGroup.POST("/confirm", handler.ConfirmOrder)    // 确认收货
	}
}
