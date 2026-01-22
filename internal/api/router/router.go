package router

import (
	adminRouter "xiaomi-mall/internal/api/router/admin"
	userRouter "xiaomi-mall/internal/api/router/user"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	// r.Use(middleware.Cors())   // 跨域
	// r.Use(middleware.Logger()) // 日志

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1 路由组
	v1 := r.Group("/api")
	{
		adminRouter.ProductRoutes(v1) // 管理员商品路由
		adminRouter.SeckillRoutes(v1) // 管理员秒杀路由

		userRouter.AddressRoutes(v1) // 用户地址路由
		userRouter.OrderRoutes(v1)   // 用户订单路由
		userRouter.ProductRoutes(v1) // 用户商品路由
		userRouter.SeckillRoutes(v1) // 用户秒杀路由
		userRouter.UserRoutes(v1)    // 用户路由
	}

	return r
}
