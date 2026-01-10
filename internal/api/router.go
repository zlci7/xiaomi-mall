package api

import (
	"net/http"
	"xiaomi-mall/internal/api/handler"

	"github.com/gin-gonic/gin" // 引入 v1 包
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// V1 版本接口
	v1Group := r.Group("/api")
	{
		// 用户模块
		v1Group.POST("/user/register", handler.UserRegister)
		v1Group.POST("/user/login", handler.UserLogin)
	}

	return r
}
