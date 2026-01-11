package router

import (
	"xiaomi-mall/internal/api/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/user")
	{
		// 公开接口（不需要登录）
		userGroup.POST("/register", handler.UserRegister)
		userGroup.POST("/login", handler.UserLogin)

		// 需要认证的接口
		// auth := userGroup.Group("")
		// auth.Use(middleware.JWTAuth()) // JWT 认证中间件
		// {
		// 	auth.GET("/profile", handler.GetUserProfile)
		// 	auth.PUT("/profile", handler.UpdateUserProfile)
		// 	auth.POST("/logout", handler.UserLogout)
		// }
	}
}
