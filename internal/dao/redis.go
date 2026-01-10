package dao

import (
	"context"
	"fmt"

	"xiaomi-mall/config"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Database.RedisAddr,
		Password: config.AppConfig.Database.RedisPw,
		DB:       0,
	})

	// 测试连接
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("❌ Redis连接失败: " + err.Error())
	}

	fmt.Println("✅ Redis 连接成功！") // ⬅️ 加上这行
}
