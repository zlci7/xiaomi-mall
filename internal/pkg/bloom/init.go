package bloom

import (
	"log"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	pkgBloom "xiaomi-mall/pkg/bloom"
)

// InitProductBloom 初始化商品布隆过滤器
func InitProductBloom() error {
	log.Println("🚀 开始初始化商品布隆过滤器...")

	// 1. 创建布隆过滤器
	// 预计 100 万商品，误判率 1%
	pkgBloom.ProductBloom = pkgBloom.NewBloomFilter(1000000, 0.01, dao.Rdb)

	// 2. 尝试从 Redis 加载
	err := pkgBloom.ProductBloom.LoadFromRedis("bloom:product")
	if err == nil {
		log.Println("✅ 商品布隆过滤器从 Redis 加载成功")
		return nil
	}

	log.Println("⚠️  Redis 中无缓存，开始从数据库重建...")

	// 3. 从数据库加载所有商品 ID
	var productIDs []uint
	err = dao.DB.Model(&model.Product{}).Pluck("id", &productIDs).Error
	if err != nil {
		return err
	}

	// 4. 添加到布隆过滤器
	for _, id := range productIDs {
		pkgBloom.ProductBloom.AddUint(id)
	}

	log.Printf("✅ 已添加 %d 个商品到布隆过滤器", len(productIDs))

	// 5. 保存到 Redis（下次启动直接加载）
	err = pkgBloom.ProductBloom.SaveToRedis("bloom:product")
	if err != nil {
		log.Printf("⚠️  保存到 Redis 失败: %v", err)
	}

	// 6. 打印统计信息
	stats := pkgBloom.ProductBloom.Stats()
	log.Printf("📊 布隆过滤器统计: %+v", stats)

	return nil
}

// InitSeckillBloom 初始化秒杀商品布隆过滤器
func InitSeckillBloom() error {
	log.Println("🚀 开始初始化秒杀布隆过滤器...")

	// 秒杀商品数量少，预计 10000 个，误判率 0.1%
	pkgBloom.SeckillBloom = pkgBloom.NewBloomFilter(10000, 0.001, dao.Rdb)

	// 尝试从 Redis 加载
	err := pkgBloom.SeckillBloom.LoadFromRedis("bloom:seckill")
	if err == nil {
		log.Println("✅ 秒杀布隆过滤器从 Redis 加载成功")
		return nil
	}

	log.Println("⚠️  Redis 中无缓存，开始从数据库重建...")

	// 从数据库加载所有秒杀商品 ID
	var seckillIDs []uint
	err = dao.DB.Table("seckill_products").Pluck("id", &seckillIDs).Error
	if err != nil {
		return err
	}

	// 添加到布隆过滤器
	for _, id := range seckillIDs {
		pkgBloom.SeckillBloom.AddUint(id)
	}

	log.Printf("✅ 已添加 %d 个秒杀商品到布隆过滤器", len(seckillIDs))

	// 保存到 Redis
	err = pkgBloom.SeckillBloom.SaveToRedis("bloom:seckill")
	if err != nil {
		log.Printf("⚠️  保存到 Redis 失败: %v", err)
	}

	return nil
}

// AddProductToBloom 添加新商品到布隆过滤器（管理员创建商品时调用）
func AddProductToBloom(productID uint) {
	if pkgBloom.ProductBloom == nil {
		return
	}

	pkgBloom.ProductBloom.AddUint(productID)

	// 异步保存到 Redis（不阻塞业务）
	go func() {
		err := pkgBloom.ProductBloom.SaveToRedis("bloom:product")
		if err != nil {
			log.Printf("⚠️  更新布隆过滤器失败: %v", err)
		}
	}()
}

// AddSeckillToBloom 添加秒杀商品到布隆过滤器
func AddSeckillToBloom(seckillID uint) {
	if pkgBloom.SeckillBloom == nil {
		return
	}

	pkgBloom.SeckillBloom.AddUint(seckillID)

	go func() {
		err := pkgBloom.SeckillBloom.SaveToRedis("bloom:seckill")
		if err != nil {
			log.Printf("⚠️  更新布隆过滤器失败: %v", err)
		}
	}()
}
