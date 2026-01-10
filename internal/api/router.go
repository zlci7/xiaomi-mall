package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 1. 健康检查接口 (用于运维检测服务是否存活)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 2. 这里的 v1 路由组稍后会放用户、商品接口
	v1 := r.Group("/api/v1")
	{
		v1.GET("test", func(c *gin.Context) {
			c.JSON(200, "API v1 works")
		})
	}

	return r
}
