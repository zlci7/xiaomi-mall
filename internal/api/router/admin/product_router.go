package adminRouter

import (
	adminHandler "xiaomi-mall/internal/api/handler/admin"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(rg *gin.RouterGroup) {
	adminGroup := rg.Group("/admin")
	{
		adminGroup.POST("/product", adminHandler.AdminCreateProduct)
		adminGroup.PUT("/product/stock", adminHandler.AdminUpdateProductStock)
		adminGroup.PUT("/product/on_sale", adminHandler.AdminToggleProductOnSale)
		// adminGroup.POST("/login", handler.AdminLogin)
	}
}
