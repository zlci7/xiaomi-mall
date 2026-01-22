package userRouter

import (
	userHandler "xiaomi-mall/internal/api/handler/user"
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// router/order_router.go
func OrderRoutes(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/order")
	orderGroup.Use(middleware.JWTAuth()) // JWT 认证
	{
		// 普通订单
		orderGroup.POST("/create", userHandler.CreateOrder)   // 创建订单
		orderGroup.POST("/pay", userHandler.PayOrder)         // 支付订单
		orderGroup.POST("/cancel", userHandler.CancelOrder)   // 取消订单
		orderGroup.GET("/:order_no", userHandler.OrderDetail) // 订单详情
		orderGroup.GET("/list", userHandler.GetOrderList)     // 订单列表
		orderGroup.POST("/confirm", userHandler.ConfirmOrder) // 确认收货
	}
}
