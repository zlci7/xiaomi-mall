package router

import (
	"xiaomi-mall/internal/api/handler"
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// router/address_router.go
func RegisterAddressRoutes(rg *gin.RouterGroup) {
	addressGroup := rg.Group("/address")
	addressGroup.Use(middleware.JWTAuth())
	{
		addressGroup.GET("/list", handler.GetAddressList)   // 地址列表
		addressGroup.POST("/save", handler.SaveAddress)     // 创建/编辑地址
		addressGroup.POST("/delete", handler.DeleteAddress) // 删除地址
	}
}
