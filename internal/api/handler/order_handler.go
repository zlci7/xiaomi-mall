package handler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 创建订单
func CreateOrder(c *gin.Context) {
	//0.获取用户ID
	userID := c.GetUint("user_id")
	//1.绑定请求参数
	var req dto.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := service.Order.CreateOrder(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}
