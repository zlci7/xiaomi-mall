package adminRouter

import (
	adminHandler "xiaomi-mall/internal/api/handler/admin"

	"github.com/gin-gonic/gin"
)

func SeckillRoutes(rg *gin.RouterGroup) {
	seckillGroup := rg.Group("/seckill")
	{
		seckillGroup.POST("/product", adminHandler.AdminCreateSeckillProduct)
		seckillGroup.DELETE("/product/:id", adminHandler.AdminDeleteSeckillProduct)
		seckillGroup.PUT("/product/:id", adminHandler.AdminUpdateSeckillStatus) // 0:未开始 1:进行中 2:已结束
	}
}
