package userRouter

import (
	userHandler "xiaomi-mall/internal/api/handler/user"
	"xiaomi-mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// 注册商品相关路由
func ProductRoutes(rg *gin.RouterGroup) {
	productGroup := rg.Group("/products")
	productGroup.Use(middleware.JWTAuth()) // JWT 认证中间件
	{
		// ✅ 查询商品列表（GET + Query Params + IP限流防爬虫）
		// productGroup.GET("", middleware.IPRateLimit(), userHandler.ProductList)
		productGroup.GET("", userHandler.ProductList)

		// ✅ 查询商品详情（GET + 路径参数 + IP限流）
		// productGroup.GET("/:product_id", middleware.IPRateLimit(), userHandler.ProductDetail)
		productGroup.GET("/:product_id", userHandler.ProductDetail)

		// ✅ 查询 SKU 详情（GET + 路径参数 + IP限流）
		// productGroup.GET("/skus/:sku_id", middleware.IPRateLimit(), userHandler.SkuDetail)
		productGroup.GET("/skus/:sku_id", userHandler.SkuDetail)
	}

	// ✅ 查询分类列表（独立路由组）
	categoryGroup := rg.Group("/categories")
	categoryGroup.Use(middleware.JWTAuth())
	{
		categoryGroup.GET("", userHandler.CategoryList) // GET /categories
	}
}
