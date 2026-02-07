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
		// 秒杀列表和详情（IP 限流：1秒100次）
		// seckillGroup.GET("/list", middleware.IPRateLimit(), userHandler.SeckillList)
		seckillGroup.GET("/list", userHandler.SeckillList)
		// seckillGroup.GET("/:id", middleware.IPRateLimit(), userHandler.SeckillDetail)
		seckillGroup.GET("/:id", userHandler.SeckillDetail)

		// 秒杀下单（严格限流：单用户1秒1次）
		seckillGroup.POST("/order",
			// middleware.SeckillRateLimit(), // 秒杀专用限流
			userHandler.CreateSeckillOrder,
		)
	}
}
