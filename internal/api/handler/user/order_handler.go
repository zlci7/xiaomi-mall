package userHandler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service/userService"
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
	resp, err := userService.Order.CreateOrder(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 支付订单
func PayOrder(c *gin.Context) {
	//0.获取用户ID
	userID := c.GetUint("user_id")
	//1.绑定请求参数
	var req dto.PayOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := userService.Order.PayOrder(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 取消订单
func CancelOrder(c *gin.Context) {
	//0.获取用户ID
	userID := c.GetUint("user_id")
	//1.绑定请求参数
	var req dto.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	err := userService.Order.CancelOrder(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, nil)
}

// 订单详情查询
func OrderDetail(c *gin.Context) {
	//0.获取用户ID
	//1.绑定请求参数
	var req dto.OrderDetailReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := userService.Order.GetOrderDetail(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 订单列表查询
func GetOrderList(c *gin.Context) {
	//0.获取用户ID
	userID := c.GetUint("user_id")
	//1.绑定请求参数
	var req dto.OrderListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := userService.Order.GetOrderList(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 确认收货
func ConfirmOrder(c *gin.Context) {
	//0.获取用户ID
	userID := c.GetUint("user_id")
	//1.绑定请求参数
	var req dto.ConfirmOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	err := userService.Order.ConfirmOrder(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, nil)
}
