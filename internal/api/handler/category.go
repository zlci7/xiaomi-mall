package handler

import (
	"xiaomi-mall/internal/service"
	"xiaomi-mall/pkg/response"

	"github.com/gin-gonic/gin"
)

// 商品分类列表查询
// GET /categories
func CategoryList(c *gin.Context) {
	//1.没有参数传递，直接调用Service
	resp, err := service.Category.CategoryList()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//2.返回响应
	response.Success(c, resp)
}
