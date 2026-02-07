package userService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	pkgBloom "xiaomi-mall/pkg/bloom"
	"xiaomi-mall/pkg/xerr"
)

type ProductService struct{}

var Product = new(ProductService)

var ctx = context.Background()

// 商品分页查询
func (s *ProductService) ProductList(req dto.ProductListReq) (*vo.ProductListResp, error) {
	// ========== 1️⃣ 设置默认值 ==========
	page := req.Page
	pageSize := req.PageSize
	onSale := req.OnSale

	// 默认第 1 页
	if page <= 0 {
		page = 1
	}

	// 默认每页 10 条
	if pageSize <= 0 {
		pageSize = 10
	}

	// 🔥 用户端默认只显示上架商品（如果需要查看全部，前端需要明确传 on_sale 参数）
	if onSale == nil {
		trueValue := true
		onSale = &trueValue
	}

	// ========== 2️⃣ 查询数据库 ==========
	products, total, err := dao.Product.GetProductList(
		req.CategoryID,
		req.Keyword,
		onSale, // ⬅️ 使用处理后的值
		req.SortBy,
		req.Order,
		page,
		pageSize,
	)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	// ========== 3️⃣ 转换为 VO ==========
	productVOs := make([]vo.ProductItemVO, 0, len(products))
	for _, product := range products {
		productVOs = append(productVOs, vo.ProductItemVO{
			ProductID:     product.ID,
			Name:          product.Name,
			Title:         product.Title,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           product.Num,
			ClickNum:      product.ClickNum,
			OnSale:        product.OnSale,
		})
	}

	// ========== 4️⃣ 返回响应 ==========
	resp := &vo.ProductListResp{
		List:     productVOs,
		Total:    total,
		Page:     page,     // ⬅️ 返回处理后的值
		PageSize: pageSize, // ⬅️ 返回处理后的值
	}
	return resp, nil
}

// 商品详情查询
func (s *ProductService) ProductDetail(req dto.ProductDetailReq) (*vo.ProductDetailResp, error) {
	// ========== 0️⃣ 布隆过滤器前置校验（防止缓存穿透）==========
	if pkgBloom.ProductBloom != nil {
		if !pkgBloom.ProductBloom.TestUint(req.ProductID) {
			// 布隆过滤器判断：商品一定不存在（100% 准确）
			println("🛡️  布隆过滤器拦截：商品不存在")
			return nil, xerr.NewErrCode(xerr.PRODUCT_NOT_FOUND)
		}
	}

	// ========== 1️⃣ 尝试从缓存读取 ==========
	cacheKey := fmt.Sprintf("product:detail:%d", req.ProductID)
	cacheData, err := dao.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 反序列化缓存数据
		var resp vo.ProductDetailResp
		if err := json.Unmarshal([]byte(cacheData), &resp); err == nil {
			println("✅ 商品详情：命中缓存") // 调试日志
			return &resp, nil
		}
		// JSON 解析失败，删除损坏的缓存
		dao.Rdb.Del(ctx, cacheKey)
	}

	println("⚠️  商品详情：缓存未命中，查询数据库") // 调试日志

	// ========== 2️⃣ 查询商品基本信息（SPU） ==========
	product, err := dao.Product.GetProductByID(req.ProductID)
	if err != nil {
		// 商品不存在，缓存空值防止缓存穿透
		dao.Rdb.Set(ctx, cacheKey, "null", 5*time.Minute)
		return nil, xerr.NewErrCode(xerr.PRODUCT_NOT_FOUND)
	}

	// ========== 3️⃣ 查询分类名称 ==========
	category, err := dao.Category.GetCategoryByID(product.CategoryID)
	if err != nil {
		// 分类不存在不影响商品展示，使用默认值
		category = &model.Category{Name: "未分类"}
	}

	// ========== 4️⃣ 查询该商品的所有 SKU ==========
	skus, err := dao.Product.GetSkusByProductID(req.ProductID)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	// 转换 SKU 为 VO（确保非 nil）
	skuVOs := make([]vo.SkuVO, 0, len(skus))
	for _, sku := range skus {
		skuVOs = append(skuVOs, vo.SkuVO{
			SkuID: sku.ID,
			Title: sku.Title,
			Price: sku.Price,
			Stock: sku.Stock,
			Code:  sku.Code,
		})
	}

	// ========== 5️⃣ 增加商品点击量（异步处理，不影响查询性能） ==========
	go func() {
		dao.Product.IncrementClickNum(req.ProductID)
	}()

	// ========== 6️⃣ 构造响应 VO ==========
	resp := &vo.ProductDetailResp{
		ProductID:     product.ID,
		Name:          product.Name,
		CategoryID:    product.CategoryID,
		CategoryName:  category.Name,
		Title:         product.Title,
		Info:          product.Info,
		ImgPath:       product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		Num:           product.Num,
		ClickNum:      product.ClickNum,
		OnSale:        product.OnSale,
		SKUs:          skuVOs, // ⬅️ 确保是 [] 而不是 null
	}

	// ========== 7️⃣ 写入缓存 ==========
	if jsonData, err := json.Marshal(resp); err == nil {
		dao.Rdb.Set(ctx, cacheKey, jsonData, time.Hour)
	}

	return resp, nil
}

// SKU详情查询
func (s *ProductService) SkuDetail(req dto.SkuDetailReq) (*vo.SkuDetailResp, error) {
	// ========== 1️⃣ 尝试从缓存读取 ==========
	cacheKey := fmt.Sprintf("sku:detail:%d", req.SkuID)
	cacheData, err := dao.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {

		// 反序列化缓存数据
		var resp vo.SkuDetailResp
		if err := json.Unmarshal([]byte(cacheData), &resp); err == nil {
			println("✅ SKU详情：命中缓存") // 调试日志
			return &resp, nil
		}
		// JSON 解析失败，删除损坏的缓存
		dao.Rdb.Del(ctx, cacheKey)
	}

	println("⚠️  SKU详情：缓存未命中，查询数据库") // 调试日志

	// ========== 2️⃣ 查询 SKU 信息 ==========
	sku, err := dao.Product.GetSkuByID(req.SkuID)
	if err != nil {
		// SKU 不存在，缓存空值防止缓存穿透
		dao.Rdb.Set(ctx, cacheKey, "null", 5*time.Minute)
		return nil, xerr.NewErrCode(xerr.PRODUCT_SKU_NOT_FOUND)
	}

	// ========== 3️⃣ 构造响应 VO ==========
	resp := &vo.SkuDetailResp{
		SkuID: sku.ID,
		Title: sku.Title,
		Price: sku.Price,
		Stock: sku.Stock,
		Code:  sku.Code,
	}

	// ========== 4️⃣ 写入缓存 ==========
	if jsonData, err := json.Marshal(resp); err == nil {
		// SKU 库存变化频繁，缓存时间设置短一些
		dao.Rdb.Set(ctx, cacheKey, jsonData, 10*time.Minute)
	}

	return resp, nil
}
