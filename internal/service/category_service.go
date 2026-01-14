package service

import (
	"context"
	"encoding/json"
	"time"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/pkg/xerr"
)

type CategoryService struct{}

var Category = new(CategoryService)

// 商品分类列表查询（带缓存）
func (s *CategoryService) CategoryList() (*vo.CategoryListResp, error) {
	ctx := context.Background()
	cacheKey := "category:list"

	// ========== 1️⃣ 尝试从 Redis 读取缓存 ==========
	cacheData, err := dao.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中！反序列化 JSON
		var resp vo.CategoryListResp
		if err := json.Unmarshal([]byte(cacheData), &resp); err == nil {
			println("✅ 分类列表：命中缓存") // ⬅️ 调试日志
			return &resp, nil      // 直接返回缓存结果
		}
		// JSON 解析失败，删除错误缓存，继续查库
		dao.Rdb.Del(ctx, cacheKey)
	}
	// err == redis.Nil 表示缓存未命中，继续查数据库
	// 其他错误（如 Redis 连接失败）也继续查库，降级处理

	println("⚠️  分类列表：缓存未命中，查询数据库") // ⬅️ 调试日志

	// ========== 2️⃣ 缓存未命中，查询数据库 ==========
	categories, err := dao.Category.GetAllCategories()
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	// ========== 3️⃣ 转换为 VO ==========
	categoryVOs := make([]vo.CategoryVO, 0, len(categories))
	for _, category := range categories {
		categoryVOs = append(categoryVOs, vo.CategoryVO{
			CategoryID:   category.ID,
			CategoryName: category.Name,
		})
	}
	resp := &vo.CategoryListResp{List: categoryVOs}

	// ========== 4️⃣ 写入 Redis 缓存 ==========
	// 序列化为 JSON
	jsonData, err := json.Marshal(resp)
	if err == nil {
		// 设置缓存，过期时间 24 小时
		dao.Rdb.Set(ctx, cacheKey, jsonData, 24*time.Hour)
		// 注意：缓存写入失败不影响业务，所以不处理错误
	}

	return resp, nil
}

// 删除分类缓存（管理员修改分类时调用）
func (s *CategoryService) DeleteCategoryCache() error {
	ctx := context.Background()
	return dao.Rdb.Del(ctx, "category:list").Err()
}
