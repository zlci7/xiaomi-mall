package userHandler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service/userService"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 秒杀商品列表查询
func SeckillList(c *gin.Context) {
	//1.绑定请求参数（URI 路径参数）
	var req dto.UserSeckillListReq
	if err := c.ShouldBindQuery(&req); err != nil { // ⬅️ 改为 ShouldBindQuery
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := userService.Seckill.GetSeckillProductList(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 秒杀商品详情查询
func SeckillDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	//1.绑定请求参数（URI 路径参数）
	var req dto.UserSeckillDetailReq
	if err := c.ShouldBindUri(&req); err != nil { // ⬅️ 改为 ShouldBindQuery
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := userService.Seckill.GetSeckillProductDetail(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}
