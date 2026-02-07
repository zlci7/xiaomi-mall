package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	rdb    *redis.Client
	ctx    context.Context
	limit  int           // 限流次数
	window time.Duration // 时间窗口
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
// limit: 时间窗口内允许的最大请求次数
// window: 时间窗口大小（如 1 秒）
func NewSlidingWindowLimiter(rdb *redis.Client, limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		rdb:    rdb,
		ctx:    context.Background(),
		limit:  limit,
		window: window,
	}
}

// Allow 检查是否允许请求
// key: 限流维度的唯一标识（如 "user:123", "ip:192.168.1.1"）
// memberID: 请求的唯一标识（使用雪花算法生成）
// 返回：是否允许，当前请求数
func (l *SlidingWindowLimiter) Allow(key string, memberID string) (bool, int, error) {
	now := time.Now().UnixMilli() // 当前时间戳（毫秒）
	windowStart := now - l.window.Milliseconds()

	// 使用 Redis Pipeline 提高性能（一次网络 IO）
	pipe := l.rdb.Pipeline()

	// 1. 删除窗口外的过期数据
	pipe.ZRemRangeByScore(l.ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// 2. 统计当前窗口内的请求数
	countCmd := pipe.ZCard(l.ctx, key)

	// 3. 执行 Pipeline
	_, err := pipe.Exec(l.ctx)
	if err != nil {
		return false, 0, err
	}

	// 获取当前请求数
	count := int(countCmd.Val())

	// 4. 判断是否超限
	if count >= l.limit {
		return false, count, nil // 拒绝请求
	}

	// 5. 添加当前请求到 ZSet
	err = l.rdb.ZAdd(l.ctx, key, &redis.Z{
		Score:  float64(now),
		Member: memberID,
	}).Err()
	if err != nil {
		return false, count, err
	}

	// 6. 设置 key 过期时间（防止内存泄漏）
	// 过期时间设为窗口的 2 倍，确保数据清理
	l.rdb.Expire(l.ctx, key, l.window*2)

	return true, count + 1, nil // 允许请求
}

// AllowN 检查是否允许 N 次请求（批量操作）
func (l *SlidingWindowLimiter) AllowN(key string, memberIDs []string) (bool, int, error) {
	now := time.Now().UnixMilli()
	windowStart := now - l.window.Milliseconds()

	pipe := l.rdb.Pipeline()

	// 1. 删除过期数据
	pipe.ZRemRangeByScore(l.ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// 2. 统计当前请求数
	countCmd := pipe.ZCard(l.ctx, key)

	_, err := pipe.Exec(l.ctx)
	if err != nil {
		return false, 0, err
	}

	count := int(countCmd.Val())

	// 3. 判断是否超限（加上即将添加的 N 个请求）
	if count+len(memberIDs) > l.limit {
		return false, count, nil
	}

	// 4. 批量添加请求
	members := make([]*redis.Z, len(memberIDs))
	for i, id := range memberIDs {
		members[i] = &redis.Z{
			Score:  float64(now),
			Member: id,
		}
	}
	err = l.rdb.ZAdd(l.ctx, key, members...).Err()
	if err != nil {
		return false, count, err
	}

	// 5. 设置过期时间
	l.rdb.Expire(l.ctx, key, l.window*2)

	return true, count + len(memberIDs), nil
}

// Reset 重置限流计数（用于测试或手动清空）
func (l *SlidingWindowLimiter) Reset(key string) error {
	return l.rdb.Del(l.ctx, key).Err()
}

// GetCount 获取当前窗口内的请求数（用于监控）
func (l *SlidingWindowLimiter) GetCount(key string) (int, error) {
	now := time.Now().UnixMilli()
	windowStart := now - l.window.Milliseconds()

	// 先清理过期数据
	l.rdb.ZRemRangeByScore(l.ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// 统计请求数
	count, err := l.rdb.ZCard(l.ctx, key).Result()
	return int(count), err
}
