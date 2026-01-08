package main

import (
	"fmt"
	"log"
	"time"

	"xiaomi-mall/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// ========== 1. 连接数据库 ==========
	dsn := "root:1234@tcp(127.0.0.1:13306)/xiaomi_mall?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印所有 SQL（开发环境）
	})

	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ 获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	fmt.Println("✅ MySQL 连接成功！")

	// ========== 2. 执行数据库迁移 ==========
	if err := model.Migrate(db); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	fmt.Println("✅ 数据库迁移完成！")

	// ========== 3. 初始化 Gin 框架 ==========
	r := gin.Default() // 自带 Logger 和 Recovery 中间件

	// ========== 4. 注册路由 ==========
	// 健康检查接口（阶段一产物要求）
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关
		v1.POST("/user/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "注册接口（待实现）"})
		})
		v1.POST("/user/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "登录接口（待实现）"})
		})

		// 商品相关
		v1.GET("/products", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "商品列表（待实现）"})
		})
	}

	// ========== 5. 启动服务 ==========
	// 监听所有网络接口，允许外部访问
	addr := "0.0.0.0:8080"
	fmt.Printf("🚀 服务启动成功，监听地址：%s\n", addr)
	fmt.Println("📌 本地访问: http://localhost:8080/ping")
	fmt.Println("📌 外部访问: http://192.168.100.128:8080/ping")

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ 启动服务失败: %v", err)
	}
}
