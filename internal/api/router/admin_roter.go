package router

import (
	"xiaomi-mall/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(rg *gin.RouterGroup) {
	adminGroup := rg.Group("/admin")
	{
		adminGroup.POST("/product", handler.AdminCreateProduct)
		adminGroup.PUT("/product/stock", handler.AdminUpdateProductStock)
		adminGroup.PUT("/product/on_sale", handler.AdminToggleProductOnSale)
		// adminGroup.POST("/login", handler.AdminLogin)
	}
}
