package userRouter

import (
	userHandler "xiaomi-mall/internal/api/handler/user"
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// router/address_router.go
func AddressRoutes(rg *gin.RouterGroup) {
	addressGroup := rg.Group("/address")
	addressGroup.Use(middleware.JWTAuth())
	{
		addressGroup.GET("/list", userHandler.GetAddressList)   // 地址列表
		addressGroup.POST("/save", userHandler.SaveAddress)     // 创建/编辑地址
		addressGroup.POST("/delete", userHandler.DeleteAddress) // 删除地址
	}
}
