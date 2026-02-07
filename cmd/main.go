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
	"xiaomi-mall/internal/pkg/bloom"
	"xiaomi-mall/internal/pkg/consumer"
	"xiaomi-mall/pkg/idgen"
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

	// 4. 初始化雪花算法（生成订单号）
	if err := idgen.InitSnowflake(1); err != nil {
		log.Fatalf("❌ 初始化雪花算法失败: %v", err)
	}
	fmt.Println("✅ 雪花算法初始化成功！")

	// 4.5 初始化布隆过滤器（防止缓存穿透）
	if err := bloom.InitProductBloom(); err != nil {
		log.Printf("⚠️  初始化商品布隆过滤器失败: %v", err)
	}
	if err := bloom.InitSeckillBloom(); err != nil {
		log.Printf("⚠️  初始化秒杀布隆过滤器失败: %v", err)
	}

	// 5. 启动秒杀订单消费者（异步写入MySQL）
	go consumer.ConsumeSeckillOrders()
	fmt.Println("✅ 秒杀订单消费者已启动")

	// 6. 启动订单超时扫描器（统一处理普通订单和秒杀订单）
	consumer.StartSeckillOrderTimeoutScanner()
	fmt.Println("✅ 订单超时扫描器已启动")

	// 7. 初始化 Gin 框架
	r := router.InitRouter()

	// 8. 启动服务（非阻塞）
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

	// 9. 优雅关闭
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
