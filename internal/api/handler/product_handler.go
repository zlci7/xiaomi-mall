package handler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 分页查询商品
// GET /products?page=1&page_size=10&category_id=10&keyword=小米手机&sort_by=price&order=desc
func ProductList(c *gin.Context) {
	//1.绑定请求参数（Query Params）
	var req dto.ProductListReq
	if err := c.ShouldBindQuery(&req); err != nil { // ⬅️ 改为 ShouldBindQuery
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := service.Product.ProductList(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 商品详情查询
// GET /products/:product_id
func ProductDetail(c *gin.Context) {
	//1.绑定请求参数（URI 路径参数）
	var req dto.ProductDetailReq
	if err := c.ShouldBindUri(&req); err != nil { // ⬅️ 改为 ShouldBindUri
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := service.Product.ProductDetail(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// SKU详情查询
// GET /products/skus/:sku_id
func SkuDetail(c *gin.Context) {
	//1.绑定请求参数（URI 路径参数）
	var req dto.SkuDetailReq
	if err := c.ShouldBindUri(&req); err != nil { // ⬅️ 改为 ShouldBindUri
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := service.Product.SkuDetail(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 商品分类列表查询
// GET /categories
func CategoryList(c *gin.Context) {
	//1.没有参数传递，直接调用Service
	resp, err := service.Product.CategoryList()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//2.返回响应
	response.Success(c, resp)
}
