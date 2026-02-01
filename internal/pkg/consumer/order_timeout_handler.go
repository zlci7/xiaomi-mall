package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/service/userService"

	"github.com/go-redis/redis/v8"
)

// StartSeckillOrderTimeoutScanner 启动秒杀订单超时扫描器
func StartSeckillOrderTimeoutScanner() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // 每秒扫描一次
		defer ticker.Stop()

		log.Println("✅ 秒杀订单超时扫描器启动")

		for range ticker.C {
			scanExpiredSeckillOrders()
		}
	}()
}

// scanExpiredSeckillOrders 扫描并处理过期的秒杀订单
func scanExpiredSeckillOrders() {
	ctx := context.Background()
	now := float64(time.Now().Unix())

	// 1. 查询已过期的订单（score <= 当前时间）
	orders, err := dao.Rdb.ZRangeByScore(ctx, "order:delay:queue", &redis.ZRangeBy{
		Min: "-inf",                 // 最小值
		Max: fmt.Sprintf("%f", now), // 当前时间
	}).Result()

	if err != nil || len(orders) == 0 {
		return // 无过期订单或查询失败
	}

	// 2. 逐个处理过期订单
	for _, orderNum := range orders {
		// 2.1 判断订单类型
		orderType, err := getOrderType(orderNum)
		if err != nil {
			// 订单不存在，从队列中移除
			dao.Rdb.ZRem(ctx, "order:delay:queue", orderNum)
			continue
		}

		// 2.2 根据订单类型选择关单方法
		var closeErr error
		if orderType == 2 { // 秒杀订单
			log.Printf("⏰ 发现过期秒杀订单：%s", orderNum)
			closeErr = userService.Seckill.CloseSeckillOrder(orderNum)
		} else { // 普通订单
			log.Printf("⏰ 发现过期普通订单：%s", orderNum)
			closeErr = userService.Order.CloseOrder(orderNum)
		}

		// 2.3 处理结果
		if closeErr == nil {
			// 关单成功，从队列中移除
			dao.Rdb.ZRem(ctx, "order:delay:queue", orderNum)
			log.Printf("✅ 订单关闭成功：%s", orderNum)
		} else {
			log.Printf("❌ 订单关闭失败：%s, 错误：%v", orderNum, closeErr)
			// 关单失败，下次继续尝试
		}
	}
}

// getOrderType 获取订单类型（1=普通订单，2=秒杀订单）
func getOrderType(orderNum string) (int, error) {
	var orderType int
	err := dao.DB.Table("orders").
		Select("type").
		Where("order_num = ?", orderNum).
		Scan(&orderType).Error
	return orderType, err
}
