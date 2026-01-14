package router

import (
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
		RegisterUserRoutes(v1)    // 用户路由
		RegisterAdminRoutes(v1)   // 管理员路由
		RegisterProductRoutes(v1) // 商品路由（将来添加）
		// RegisterOrderRoutes(v1)   // 订单路由（将来添加）
	}

	return r
}
