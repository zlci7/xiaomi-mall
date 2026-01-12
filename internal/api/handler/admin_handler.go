package handler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 最小实现：只做必要的 3 个接口
func AdminCreateProduct(c *gin.Context) {
	//1.绑定请求参数
	var req dto.CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := service.Product.CreateProduct(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)

}

func AdminUpdateProductStock(c *gin.Context) {
	// 补充库存（运营常用）
}

func AdminToggleProductOnSale(c *gin.Context) {
	// 上架/下架（简单的状态切换）
}
