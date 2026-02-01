package userRouter

import (
	userHandler "xiaomi-mall/internal/api/handler/user"
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SeckillRoutes(rg *gin.RouterGroup) {
	seckillGroup := rg.Group("/seckill")
	seckillGroup.Use(middleware.JWTAuth()) // JWT 认证
	{
		seckillGroup.GET("/list", userHandler.SeckillList)  // 活动列表
		seckillGroup.GET("/:id", userHandler.SeckillDetail) // 活动详情
		// // 秒杀下单（需要登录）
		seckillGroup.POST("/order", userHandler.CreateSeckillOrder) // 秒杀下单
	}
}
