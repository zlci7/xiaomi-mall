package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"xiaomi-mall/config"
	"xiaomi-mall/internal/api/router"
	"xiaomi-mall/internal/dao"
)

func main() {
	// 1. 初始化配置（在所有操作之前）
	if err := config.InitConfig("../config"); err != nil { // ⬅️ 修改路径
		log.Fatalf("❌ 初始化配置失败: %v", err)
	}
	fmt.Println("✅ 配置加载成功！")

	// 2. 初始化数据库 (MySQL)
	dao.InitMySQL()

	// 3. 初始化 Redis
	dao.InitRedis()

	// 4. 初始化 Gin 框架
	r := router.InitRouter()

	// 5. 启动服务（非阻塞）
	addr := config.AppConfig.Server.Port
	fmt.Printf("🚀 服务启动成功，监听地址：%s\n", addr)
	fmt.Println("📌 本地访问: http://localhost" + addr + "/ping")
	fmt.Println("📌 外部访问: http://192.168.100.128" + addr + "/ping")

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("❌ 启动服务失败: %v\n", err)
		}
	}()

	// 临时调试
	fmt.Println("JWT Secret Key:", config.AppConfig.Jwt.AccessSecret)
	fmt.Println("JWT Expire:", config.AppConfig.Jwt.AccessExpire)

	// 6. 关闭服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 收到停止信号，正在关闭...")

	// 关闭数据库连接
	if sqlDB, err := dao.DB.DB(); err == nil {
		sqlDB.Close()
		fmt.Println("✅ 数据库连接已关闭")
	}

	// 关闭 Redis 连接
	if dao.Rdb != nil {
		dao.Rdb.Close()
		fmt.Println("✅ Redis 连接已关闭")
	}

	fmt.Println("✅ 服务已关闭")
}
