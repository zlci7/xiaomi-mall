package adminHandler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service/adminService"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 创建秒杀商品
func AdminCreateSeckillProduct(c *gin.Context) {
	//1.绑定请求参数
	var req dto.CreateSeckillProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	resp, err := adminService.Seckill.CreateSeckillProduct(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, resp)
}

// 删除秒杀商品
func AdminDeleteSeckillProduct(c *gin.Context) {
	//1.绑定请求参数
	var req dto.DeleteSeckillProductReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	err := adminService.Seckill.DeleteSeckillProduct(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, nil)
}

// 开始/结束秒杀
func AdminUpdateSeckillStatus(c *gin.Context) {
	//1.绑定请求参数
	var req dto.UpdateSeckillStatusReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "URI绑定失败: "+err.Error())
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "Query绑定失败: "+err.Error())
		return
	}
	//2.调用Service
	err := adminService.Seckill.UpdateSeckillStatus(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, nil)
}

// 预热秒杀商品到 Redis
func AdminPreheatSeckillProduct(c *gin.Context) {
	//1.绑定请求参数
	var req dto.PreheatSeckillProductReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}
	//2.调用Service
	err := adminService.Seckill.PreheatSeckillProduct(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	//3.返回响应
	response.Success(c, nil)
}
