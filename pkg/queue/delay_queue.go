package queue

import (
	"context"
	"fmt"
	"time"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/service"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// StartDelayQueueScanner 启动延迟队列扫描器（后台协程）
func StartDelayQueueScanner() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // 每秒扫描一次
		defer ticker.Stop()

		fmt.Println("✅ 延迟队列扫描器启动")

		for range ticker.C {
			ScanExpiredOrders()
		}
	}()
}

// ScanExpiredOrders 扫描并处理过期订单
func ScanExpiredOrders() {
	now := float64(time.Now().Unix())

	// 1️⃣ 查询已过期的订单（score <= 当前时间）
	orders, err := dao.Rdb.ZRangeByScore(ctx, "order:delay:queue", &redis.ZRangeBy{
		Min: "-inf",                 // 最小值
		Max: fmt.Sprintf("%f", now), // 当前时间
	}).Result()

	if err != nil || len(orders) == 0 {
		return // 无过期订单或查询失败
	}

	// 2️⃣ 逐个处理过期订单
	for _, orderNo := range orders {
		fmt.Printf("⏰ 发现过期订单：%s\n", orderNo)

		// 执行关单逻辑
		if err := service.Order.CloseOrder(orderNo); err == nil {
			// 3️⃣ 关单成功，从队列中移除
			dao.Rdb.ZRem(ctx, "order:delay:queue", orderNo)
			fmt.Printf("✅ 订单关闭成功：%s\n", orderNo)
		} else {
			fmt.Printf("❌ 订单关闭失败：%s, 错误：%v\n", orderNo, err)
			// 关单失败，下次继续尝试
		}
	}
}
