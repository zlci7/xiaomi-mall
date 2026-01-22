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
		// ✅ 查询商品列表（GET + Query Params）
		productGroup.GET("", userHandler.ProductList) // GET /products?page=1&category_id=10

		// ✅ 查询商品详情（GET + 路径参数）
		productGroup.GET("/:product_id", userHandler.ProductDetail) // GET /products/123

		// ✅ 查询 SKU 详情（GET + 路径参数）
		productGroup.GET("/skus/:sku_id", userHandler.SkuDetail) // GET /products/skus/456
	}

	// ✅ 查询分类列表（独立路由组）
	categoryGroup := rg.Group("/categories")
	categoryGroup.Use(middleware.JWTAuth())
	{
		categoryGroup.GET("", userHandler.CategoryList) // GET /categories
	}
}
