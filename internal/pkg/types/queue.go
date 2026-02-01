package types

// SeckillOrderQueueData 秒杀订单队列数据
// 用于 Redis 队列和延迟队列的数据传递
type SeckillOrderQueueData struct {
	UserID       uint   `json:"user_id"`
	SeckillID    uint   `json:"seckill_id"`
	OrderNum     string `json:"order_num"`
	Timestamp    int64  `json:"timestamp"`
	RetryCount   int    `json:"retry_count"`    // 重试次数
	FirstTryTime int64  `json:"first_try_time"` // 首次尝试时间
	LastTryTime  int64  `json:"last_try_time"`  // 最后尝试时间
}

// DeadLetterData 死信队列数据
type DeadLetterData struct {
	OrderData *SeckillOrderQueueData `json:"order_data"`
	LastError string                 `json:"last_error"`
	FailedAt  int64                  `json:"failed_at"`
}
