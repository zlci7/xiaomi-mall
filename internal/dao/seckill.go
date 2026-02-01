package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"xiaomi-mall/internal/model"

	"github.com/go-redis/redis/v8"
)

type SeckillDao struct{}

var Seckill = new(SeckillDao)

// ==================== 管理端：秒杀商品管理 ====================

// 创建秒杀商品入库
func (d *SeckillDao) CreateSeckillProduct(product *model.SeckillProduct) error {
	return DB.Create(product).Error

}

// 删除秒杀商品
func (d *SeckillDao) DeleteSeckillProduct(id uint) error {
	// return DB.Model(&model.SeckillProduct{}).Where("id =?", id).Delete(&model.SeckillProduct{}).Error
	return DB.Delete(&model.SeckillProduct{}, "id = ?", id).Error
}

// 手动开启/结束秒杀
func (d *SeckillDao) UpdateSeckillStatus(id uint, status int) error {
	return DB.Model(&model.SeckillProduct{}).Where("id=?", id).Update("status", status).Error
}

// 查询秒杀商品详情(用于管理端查询和预热)
func (d *SeckillDao) GetSeckillProductByID(id uint) (*model.SeckillProduct, error) {
	var seckillProduct model.SeckillProduct
	err := DB.Model(&model.SeckillProduct{}).Where("id = ?", id).First(&seckillProduct).Error
	return &seckillProduct, err
}

// ==================== 用户端：秒杀商品查询 ====================

// ==================== 用户端：秒杀商品下单 ====================

// ==================== 管理端：redis管理 ====================

// 设置redis库存
// ==================== Redis 操作 ====================

// 设置秒杀库存
func (d *SeckillDao) SetSeckillStock(ctx context.Context, seckillID uint, stock uint, ttl time.Duration) error {
	key := fmt.Sprintf("seckill:stock:%d", seckillID)
	return Rdb.SetNX(ctx, key, stock, ttl).Err()
}

// 获取秒杀库存
func (d *SeckillDao) GetSeckillStock(ctx context.Context, seckillID uint) (int, error) {
	key := fmt.Sprintf("seckill:stock:%d", seckillID)
	return Rdb.Get(ctx, key).Int()
}

// 扣减秒杀库存
func (d *SeckillDao) DecrSeckillStock(ctx context.Context, seckillID uint) (int64, error) {
	key := fmt.Sprintf("seckill:stock:%d", seckillID)
	return Rdb.Decr(ctx, key).Result()
}

// 缓存秒杀商品详情
func (d *SeckillDao) SetSeckillProductCache(ctx context.Context, seckillID uint, data interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("seckill:product:%d", seckillID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Rdb.Set(ctx, key, jsonData, ttl).Err()
}

// 获取秒杀商品详情
func (d *SeckillDao) GetSeckillProductCache(ctx context.Context, seckillID uint) (string, error) {
	key := fmt.Sprintf("seckill:product:%d", seckillID)
	return Rdb.Get(ctx, key).Result()
}

// 删除秒杀缓存（用于回滚）
func (d *SeckillDao) DeleteSeckillCache(ctx context.Context, seckillID uint) error {
	stockKey := fmt.Sprintf("seckill:stock:%d", seckillID)
	productKey := fmt.Sprintf("seckill:product:%d", seckillID)
	return Rdb.Del(ctx, stockKey, productKey).Err()
}

// 设置用户购买标记
func (d *SeckillDao) SetUserPurchased(ctx context.Context, seckillID uint, userID uint, ttl time.Duration) error {
	key := fmt.Sprintf("seckill:user:%d:%d", seckillID, userID)
	return Rdb.Set(ctx, key, "1", ttl).Err()
}

// 原子性地设置库存和商品详情
func (d *SeckillDao) PreheatSeckillAtomic(ctx context.Context, product *model.SeckillProduct, cacheData []byte, ttl int64) error {
	// Lua 脚本：原子性地设置库存和商品详情
	script := `
		-- 原子性预热
		local stock_key = KEYS[1]
		local product_key = KEYS[2]
		local start_list = KEYS[3]
		local end_list = KEYS[4]

		local stock = ARGV[1]
		local product_data = ARGV[2]
		local ttl = ARGV[3]
		local seckill_id = ARGV[4]
		local start_time = ARGV[5]
		local end_time = ARGV[6]

		-- 检查是否已预热
		if redis.call('EXISTS', stock_key) == 1 then
			return 0
		end

		-- 1. 设置库存
		redis.call('SET', stock_key, stock, 'EX', ttl)

		-- 2. 设置商品详情
		redis.call('SET', product_key, product_data, 'EX', ttl)

		-- 3. 添加到两个列表
		redis.call('ZADD', start_list, start_time, seckill_id)
		redis.call('ZADD', end_list, end_time, seckill_id)

		return 1
    `

	stockKey := fmt.Sprintf("seckill:stock:%d", product.ID)
	productKey := fmt.Sprintf("seckill:product:%d", product.ID)
	startList := "seckill:active:start" // 全局开始时间列表
	endList := "seckill:active:end"     // 全局结束时间列表
	result, err := Rdb.Eval(ctx, script,
		[]string{stockKey, productKey, startList, endList},
		product.SeckillStock,
		cacheData,
		ttl,
		product.ID,
		product.StartTime.Unix(),
		product.EndTime.Unix(),
	).Int()

	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("库存已预热")
	}
	return nil
}

// ==================== 用户端：秒杀查询（Redis） ====================

// 1. 获取未结束的秒杀ID列表（分页）
func (d *SeckillDao) GetActiveSeckillIDs(ctx context.Context, page, pageSize int) ([]uint, int64, error) {
	now := time.Now().Unix()

	// 从结束时间列表查询：结束时间 > 当前时间（未结束的）
	members, err := Rdb.ZRangeByScore(ctx, "seckill:active:end", &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", now),
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, 0, err
	}

	total := int64(len(members))

	// 应用层分页
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(members) {
		return []uint{}, total, nil
	}
	if end > len(members) {
		end = len(members)
	}

	pageMembers := members[start:end]

	// 转换为 uint
	ids := make([]uint, 0, len(pageMembers))
	for _, member := range pageMembers {
		id, err := strconv.ParseUint(member, 10, 64)
		if err == nil {
			ids = append(ids, uint(id))
		}
	}

	return ids, total, nil
}

// 2. 批量获取秒杀商品详情（Pipeline）
func (d *SeckillDao) BatchGetSeckillProductCache(ctx context.Context, ids []uint) (map[uint][]byte, error) {
	if len(ids) == 0 {
		return map[uint][]byte{}, nil
	}

	pipe := Rdb.Pipeline()
	cmds := make(map[uint]*redis.StringCmd)

	for _, id := range ids {
		key := fmt.Sprintf("seckill:product:%d", id)
		cmds[id] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[uint][]byte)
	for id, cmd := range cmds {
		data, err := cmd.Bytes()
		if err == nil {
			result[id] = data
		}
	}

	return result, nil
}

// 3. 批量获取库存（Pipeline）
func (d *SeckillDao) BatchGetSeckillStocks(ctx context.Context, ids []uint) (map[uint]int, error) {
	if len(ids) == 0 {
		return map[uint]int{}, nil
	}

	pipe := Rdb.Pipeline()
	cmds := make(map[uint]*redis.StringCmd)

	for _, id := range ids {
		key := fmt.Sprintf("seckill:stock:%d", id)
		cmds[id] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[uint]int)
	for id, cmd := range cmds {
		stock, err := cmd.Int()
		if err == nil {
			result[id] = stock
		} else {
			result[id] = 0
		}
	}

	return result, nil
}

// 4. 获取单个秒杀商品详情
func (d *SeckillDao) GetSeckillProductCacheByID(ctx context.Context, id uint) ([]byte, error) {
	key := fmt.Sprintf("seckill:product:%d", id)
	return Rdb.Get(ctx, key).Bytes()
}

// 5. 获取单个秒杀商品库存
func (d *SeckillDao) GetSeckillStockByID(ctx context.Context, id uint) (int, error) {
	key := fmt.Sprintf("seckill:stock:%d", id)
	return Rdb.Get(ctx, key).Int()
}

// 6. 检查用户是否已购买
func (d *SeckillDao) CheckUserPurchased(ctx context.Context, seckillID, userID uint) (bool, error) {
	key := fmt.Sprintf("seckill:user:%d:%d", seckillID, userID)
	exists, err := Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// 原子性地创建秒杀订单
func (d *SeckillDao) CreateSeckillOrderAtomic(ctx context.Context, userID, seckillID uint, orderNum string) (int, error) {
	script := `
		-- 原子性创建秒杀订单
		local stock_key = KEYS[1]
		local user_key = KEYS[2]
		local order_queue_key = KEYS[3]

		local user_id = ARGV[1]
		local seckill_id = ARGV[2]
		local order_num = ARGV[3]
		local ttl = ARGV[4]

		-- 1. 检查用户是否购买
		if redis.call('EXISTS', user_key) == 1 then
			return -1
		end

		-- 2. 检查库存
		if redis.call('EXISTS', stock_key) == 0 then
			return -3
		end

		local stock = tonumber(redis.call('GET', stock_key))
		if not stock or stock <= 0 then
			return -2
		end

		-- 3. 扣减库存
		redis.call('DECR', stock_key)

		-- 4. 设置用户购买标记
		redis.call('SETEX', user_key, ttl, order_num)

		-- 5. 将订单信息加入到队列（异步处理）
		local time_result = redis.call('TIME')
		local timestamp = tonumber(time_result[1])
		local order_data = {
			user_id = user_id,
			seckill_id = seckill_id,
			order_num = order_num,
			timestamp = timestamp,
			retry_count = 0,
			first_try_time = timestamp,
			last_try_time = timestamp,
		}
		redis.call('LPUSH', order_queue_key, cjson.encode(order_data))
		return 1
    `

	stockKey := fmt.Sprintf("seckill:stock:%d", seckillID)
	userKey := fmt.Sprintf("seckill:user:%d:%d", seckillID, userID)
	orderQueueKey := "seckill:order:queue"
	result, err := Rdb.Eval(ctx, script,
		[]string{stockKey, userKey, orderQueueKey},
		userID,
		seckillID,
		orderNum,
		86400,
	).Int()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// ==================== Redis 活动列表管理 ====================

// AddToActiveList 添加秒杀商品到活动列表
func (d *SeckillDao) AddToActiveList(ctx context.Context, seckillID uint, startTime, endTime time.Time) error {
	startList := "seckill:active:start"
	endList := "seckill:active:end"

	pipe := Rdb.Pipeline()
	pipe.ZAdd(ctx, startList, &redis.Z{
		Score:  float64(startTime.Unix()),
		Member: seckillID,
	})
	pipe.ZAdd(ctx, endList, &redis.Z{
		Score:  float64(endTime.Unix()),
		Member: seckillID,
	})

	_, err := pipe.Exec(ctx)
	return err
}

// RemoveFromActiveList 从活动列表中移除秒杀商品
func (d *SeckillDao) RemoveFromActiveList(ctx context.Context, seckillID uint) error {
	startList := "seckill:active:start"
	endList := "seckill:active:end"

	pipe := Rdb.Pipeline()
	pipe.ZRem(ctx, startList, seckillID)
	pipe.ZRem(ctx, endList, seckillID)

	_, err := pipe.Exec(ctx)
	return err
}
