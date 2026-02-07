package middleware

import (
	"fmt"
	"time"

	"xiaomi-mall/internal/dao"
	"xiaomi-mall/pkg/idgen"
	"xiaomi-mall/pkg/ratelimit"
	"xiaomi-mall/pkg/response"
	"xiaomi-mall/pkg/xerr"

	"github.com/gin-gonic/gin"
)

// 全局限流器（单例）
var (
	GlobalLimiter *ratelimit.SlidingWindowLimiter // 全局限流
	IPLimiter     *ratelimit.SlidingWindowLimiter // IP 限流
	UserLimiter   *ratelimit.SlidingWindowLimiter // 用户限流
	APILimiter    *ratelimit.SlidingWindowLimiter // 接口限流
)

// InitRateLimiters 初始化限流器
func InitRateLimiters() {
	// 全局限流：1 秒内 10000 次
	GlobalLimiter = ratelimit.NewSlidingWindowLimiter(dao.Rdb, 10000, 1*time.Second)

	// IP 限流：1 秒内 100 次
	IPLimiter = ratelimit.NewSlidingWindowLimiter(dao.Rdb, 100, 1*time.Second)

	// 用户限流：1 秒内 10 次
	UserLimiter = ratelimit.NewSlidingWindowLimiter(dao.Rdb, 10, 1*time.Second)

	// 接口限流：1 秒内 1000 次
	APILimiter = ratelimit.NewSlidingWindowLimiter(dao.Rdb, 1000, 1*time.Second)
}

// GlobalRateLimit 全局限流中间件
func GlobalRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GlobalLimiter == nil {
			c.Next()
			return
		}

		// 生成请求唯一标识
		memberID := idgen.GenStringID()

		// 检查是否允许
		allowed, count, err := GlobalLimiter.Allow("rate_limit:global", memberID)
		if err != nil {
			// Redis 错误不影响业务（降级策略）
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, xerr.RATE_LIMIT_ERROR, fmt.Sprintf("系统繁忙，请稍后重试（当前QPS: %d）", count))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimit IP 限流中间件
func IPRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IPLimiter == nil {
			c.Next()
			return
		}

		// 获取客户端 IP
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:ip:%s", ip)
		memberID := idgen.GenStringID()

		allowed, count, err := IPLimiter.Allow(key, memberID)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, xerr.RATE_LIMIT_ERROR, fmt.Sprintf("请求过于频繁，请稍后重试（当前: %d次/秒）", count))
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimit 用户限流中间件（需要先经过 JWT 认证）
func UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if UserLimiter == nil {
			c.Next()
			return
		}

		// 从 JWT 中间件获取 user_id
		userID, exists := c.Get("user_id")
		if !exists {
			// 未登录用户不限流（或可以用 IP 限流）
			c.Next()
			return
		}

		key := fmt.Sprintf("rate_limit:user:%v", userID)
		memberID := idgen.GenStringID()

		allowed, count, err := UserLimiter.Allow(key, memberID)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, xerr.RATE_LIMIT_ERROR, fmt.Sprintf("操作过于频繁，请稍后重试（当前: %d次/秒）", count))
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIRateLimit 接口限流中间件（针对特定接口）
// 使用方式：r.GET("/api/seckill", APIRateLimit("/api/seckill", 100, 1*time.Second), handler)
func APIRateLimit(apiPath string, limit int, window time.Duration) gin.HandlerFunc {
	limiter := ratelimit.NewSlidingWindowLimiter(dao.Rdb, limit, window)

	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:api:%s", apiPath)
		memberID := idgen.GenStringID()

		allowed, count, err := limiter.Allow(key, memberID)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, xerr.RATE_LIMIT_ERROR, fmt.Sprintf("接口请求过于频繁（当前: %d次/%v）", count, window))
			c.Abort()
			return
		}

		c.Next()
	}
}

// SeckillRateLimit 秒杀专用限流（更严格）
// 1 秒内单个用户最多 1 次
func SeckillRateLimit() gin.HandlerFunc {
	limiter := ratelimit.NewSlidingWindowLimiter(dao.Rdb, 1, 1*time.Second)

	return func(c *gin.Context) {
		// 必须登录
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, xerr.TOKEN_NOT_EXIST, "请先登录")
			c.Abort()
			return
		}

		key := fmt.Sprintf("rate_limit:seckill:user:%v", userID)
		memberID := idgen.GenStringID()

		allowed, _, err := limiter.Allow(key, memberID)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, xerr.RATE_LIMIT_ERROR, "操作过于频繁，请1秒后重试")
			c.Abort()
			return
		}

		c.Next()
	}
}
