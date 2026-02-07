package bloom

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis/v8"
)

// BloomFilter 布隆过滤器管理器
type BloomFilter struct {
	filter *bloom.BloomFilter
	rdb    *redis.Client
	ctx    context.Context
}

var (
	// ProductBloom 商品布隆过滤器（全局单例）
	ProductBloom *BloomFilter

	// SeckillBloom 秒杀商品布隆过滤器（全局单例）
	SeckillBloom *BloomFilter
)

// NewBloomFilter 创建布隆过滤器
// n: 预计元素数量
// p: 误判率（0.01 = 1%）
func NewBloomFilter(n uint, p float64, rdb *redis.Client) *BloomFilter {
	return &BloomFilter{
		filter: bloom.NewWithEstimates(n, p),
		rdb:    rdb,
		ctx:    context.Background(),
	}
}

// Add 添加元素
func (bf *BloomFilter) Add(item string) {
	bf.filter.AddString(item)
}

// AddUint 添加 uint 类型（商品 ID）
func (bf *BloomFilter) AddUint(id uint) {
	bf.filter.AddString(fmt.Sprintf("%d", id))
}

// Test 检查元素是否存在
// 返回 false：一定不存在（100% 准确）
// 返回 true：可能存在（99% 准确）
func (bf *BloomFilter) Test(item string) bool {
	return bf.filter.TestString(item)
}

// TestUint 检查 uint 类型
func (bf *BloomFilter) TestUint(id uint) bool {
	return bf.filter.TestString(fmt.Sprintf("%d", id))
}

// SaveToRedis 保存到 Redis（持久化）
func (bf *BloomFilter) SaveToRedis(key string) error {
	// 序列化布隆过滤器
	data, err := bf.filter.GobEncode()
	if err != nil {
		return fmt.Errorf("序列化布隆过滤器失败: %v", err)
	}

	// Base64 编码（Redis 存储更安全）
	encoded := base64.StdEncoding.EncodeToString(data)

	// 存储到 Redis（7 天过期）
	err = bf.rdb.Set(bf.ctx, key, encoded, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("保存到 Redis 失败: %v", err)
	}

	log.Printf("✅ 布隆过滤器已保存到 Redis: %s", key)
	return nil
}

// LoadFromRedis 从 Redis 加载
func (bf *BloomFilter) LoadFromRedis(key string) error {
	// 从 Redis 读取
	encoded, err := bf.rdb.Get(bf.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("布隆过滤器不存在: %s", key)
	}
	if err != nil {
		return fmt.Errorf("从 Redis 读取失败: %v", err)
	}

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("解码失败: %v", err)
	}

	// 反序列化
	err = bf.filter.GobDecode(data)
	if err != nil {
		return fmt.Errorf("反序列化失败: %v", err)
	}

	log.Printf("✅ 布隆过滤器已从 Redis 加载: %s", key)
	return nil
}

// Stats 获取统计信息
func (bf *BloomFilter) Stats() map[string]interface{} {
	return map[string]interface{}{
		"capacity":        bf.filter.Cap(),              // 位数组容量
		"hash_functions":  bf.filter.K(),                // 哈希函数数量
		"estimated_items": bf.filter.ApproximatedSize(), // 估算元素数量
	}
}
