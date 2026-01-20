package router

import (
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSeckillRoutes(rg *gin.RouterGroup) {
	seckillGroup := rg.Group("/seckill")
	seckillGroup.Use(middleware.JWTAuth()) // JWT 认证

	{
		// // 秒杀活动（无需登录可查看）
		// seckillGroup.GET("/activity/list", handler.GetSeckillActivityList)           // 活动列表
		// seckillGroup.GET("/activity/:activity_id", handler.GetSeckillActivityDetail) // 活动详情
		// // 秒杀下单（需要登录）
		// seckillGroup.POST("/order/create", handler.CreateSeckillOrder) // 秒杀下单
	}
}
