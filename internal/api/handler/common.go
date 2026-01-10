package handler

import (
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// handleServiceError 统一处理 Service 层错误（避免重复代码）
func handleServiceError(c *gin.Context, err error) {
	if codeErr, ok := err.(*xerr.CodeError); ok {
		response.Error(c, codeErr.GetErrCode(), "")
	} else {
		response.Error(c, xerr.SERVER_COMMON_ERROR, "")
	}
}
