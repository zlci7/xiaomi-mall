package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	"xiaomi-mall/internal/pkg/types"

	"gorm.io/gorm"
)

// ConsumeSeckillOrders 消费秒杀订单队列（异步写入MySQL）
func ConsumeSeckillOrders() {
	ctx := context.Background()
	queueKey := "seckill:order:queue"

	log.Println("🚀 秒杀订单消费者启动...")

	for {
		// 1. 阻塞获取订单（5秒超时）
		result, err := dao.Rdb.BRPop(ctx, 5*time.Second, queueKey).Result()
		if err != nil {
			continue // 超时或错误，继续等待
		}

		// 2. 解析订单数据
		var orderData types.SeckillOrderQueueData
		if err := json.Unmarshal([]byte(result[1]), &orderData); err != nil {
			log.Printf("❌ 订单数据解析失败: %v", err)
			continue
		}

		// 3. 写入 MySQL
		if err := writeSeckillOrderToDB(&orderData); err != nil {
			log.Printf("❌ 订单入库失败: %s, 进入重试...", orderData.OrderNum)
			handleRetry(ctx, &orderData, err)
		} else {
			log.Printf("✅ 订单入库成功: %s", orderData.OrderNum)
		}
	}
}

// 写入数据库（事务）
func writeSeckillOrderToDB(orderData *types.SeckillOrderQueueData) error {
	// 1. 查询秒杀商品信息
	seckillProduct, err := dao.Seckill.GetSeckillProductByID(orderData.SeckillID)
	if err != nil {
		return err
	}

	// 2. 查询商品信息（SPU）
	product, err := dao.Product.GetProductByID(seckillProduct.ProductID)
	if err != nil {
		return err
	}

	// 3. 查询 SKU 信息
	sku, err := dao.Product.GetSkuByID(seckillProduct.SkuID)
	if err != nil {
		return err
	}

	// 4. 事务写入
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		// 4.1 写入秒杀订单表
		seckillOrder := &model.SeckillOrder{
			UserID:           orderData.UserID,
			SeckillProductID: orderData.SeckillID,
			OrderNum:         orderData.OrderNum,
			Status:           0, // 待支付
		}
		if err := tx.Create(seckillOrder).Error; err != nil {
			return err
		}

		// 4.2 写入主订单表
		order := &model.Order{
			UserID:      orderData.UserID,
			OrderNum:    orderData.OrderNum,
			AllPrice:    int64(seckillProduct.SeckillPrice),
			PayStatus:   0,
			OrderStatus: 0,
			Type:        2, // 秒杀订单
			ExpireTime:  time.Now().Add(30 * time.Minute),
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 4.3 写入订单详情表（商品快照）
		orderItem := &model.OrderItem{
			OrderNum:     order.OrderNum,
			ProductID:    product.ID,
			ProductSkuID: sku.ID,
			Num:          1,                                  // 秒杀固定1件
			Price:        int64(seckillProduct.SeckillPrice), // 秒杀价
			Title:        product.Name + " - " + sku.Title,   // 商品名 + SKU规格
			ImgPath:      product.ImgPath,                    // 商品图片
		}
		if err := tx.Create(orderItem).Error; err != nil {
			return err
		}

		return nil
	})
}

// 重试逻辑（指数退避）
func handleRetry(ctx context.Context, orderData *types.SeckillOrderQueueData, lastErr error) {
	orderData.RetryCount++
	orderData.LastTryTime = time.Now().Unix()

	// 首次失败时记录初始时间
	if orderData.FirstTryTime == 0 {
		orderData.FirstTryTime = orderData.Timestamp
	}

	if orderData.RetryCount >= 5 {
		// 超过最大重试次数，投递到死信队列
		moveToDeadLetter(ctx, orderData, lastErr)
		return
	}

	// 计算延迟时间（指数退避：2s, 4s, 8s, 16s, 32s）
	delay := time.Duration(1<<uint(orderData.RetryCount)) * 2 * time.Second
	log.Printf("⚠️  订单入库失败，将在 %v 后重试：%s (第 %d 次重试)",
		delay, orderData.OrderNum, orderData.RetryCount)

	time.Sleep(delay)

	// 重新投递到队列
	data, _ := json.Marshal(orderData)
	dao.Rdb.LPush(ctx, "seckill:order:queue", data)
}

// 投递到死信队列
func moveToDeadLetter(ctx context.Context, orderData *types.SeckillOrderQueueData, lastErr error) {
	deadLetter := types.DeadLetterData{
		OrderData: orderData,
		LastError: lastErr.Error(),
		FailedAt:  time.Now().Unix(),
	}
	data, _ := json.Marshal(deadLetter)
	dao.Rdb.LPush(ctx, "seckill:order:dead", data)

	log.Printf("🚨 订单入库失败，已投递死信队列: %s", orderData.OrderNum)
	// TODO: 发送告警
}
