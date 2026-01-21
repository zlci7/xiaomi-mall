package handler

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/service"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// ========== 获取地址列表 ==========
// GET /api/v1/address/list
func GetAddressList(c *gin.Context) {
	// 1. 获取用户 ID（从 JWT 中间件）
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, xerr.USER_NOT_LOGIN, "")
		return
	}

	// 2. 调用 Service
	resp, err := service.Address.GetAddressList(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// 3. 返回响应
	response.Success(c, resp)
}

// ========== 保存地址（创建/编辑）==========
// POST /api/v1/address/save
func SaveAddress(c *gin.Context) {
	// 1. 获取用户 ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, xerr.USER_NOT_LOGIN, "")
		return
	}

	// 2. 绑定请求参数
	var req dto.SaveAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}

	// 3. 调用 Service
	if err := service.Address.SaveAddress(userID, req); err != nil {
		handleServiceError(c, err)
		return
	}

	// 4. 返回响应
	response.Success(c, nil)
}

// ========== 删除地址 ==========
// POST /api/v1/address/delete
func DeleteAddress(c *gin.Context) {
	// 1. 获取用户 ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, xerr.USER_NOT_LOGIN, "")
		return
	}

	// 2. 绑定请求参数
	var req dto.DeleteAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.REUQEST_PARAM_ERROR, "")
		return
	}

	// 3. 调用 Service
	if err := service.Address.DeleteAddress(userID, req); err != nil {
		handleServiceError(c, err)
		return
	}

	// 4. 返回响应
	response.Success(c, nil)
}
