package dao

import (
	"xiaomi-mall/internal/model"

	"gorm.io/gorm"
)

var Product = new(ProductDao)

type ProductDao struct{}

// ============ 商品查询（基础CRUD） ============

// 1. 商品列表查询（支持所有筛选条件）
func (d *ProductDao) GetProductList(
	categoryID uint,
	keyword string,
	onSale *bool,
	sortBy string,
	order string,
	page, pageSize int,
) (products []*model.Product, total int64, err error) {
	query := DB.Model(&model.Product{})

	// ========== 筛选条件 ==========

	// 1. 按分类筛选（可选）
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 2. 按关键词搜索（可选，模糊查询商品名）
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 3. 按上架状态筛选（可选）
	if onSale != nil {
		query = query.Where("on_sale = ?", *onSale)
	}

	// ========== 统计总数 ==========
	query.Count(&total)

	// ========== 排序 ==========
	// 默认：按 ID 倒序
	orderClause := "id DESC"

	// 自定义排序
	if sortBy != "" {
		// 防止 SQL 注入：只允许白名单字段
		validSortFields := map[string]bool{
			"price":      true, // 价格
			"num":        true, // 销量
			"click_num":  true, // 点击量
			"created_at": true, // 创建时间
			"id":         true, // ID
		}

		if validSortFields[sortBy] {
			// 排序方向（默认 DESC）
			if order == "asc" {
				orderClause = sortBy + " ASC"
			} else {
				orderClause = sortBy + " DESC"
			}
		}
	}

	// ========== 分页查询 ==========
	offset := (page - 1) * pageSize
	err = query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&products).Error
	return
}

// 2. 商品详情查询（根据商品ID）
func (d *ProductDao) GetProductByID(productID uint) (product *model.Product, err error) {
	err = DB.Model(&model.Product{}).Where("id = ?", productID).First(&product).Error
	return
}

// 3. 增加商品点击量（浏览统计）
func (d *ProductDao) IncrementClickNum(productID uint) error {
	return DB.Model(&model.Product{}).Where("id = ?", productID).
		UpdateColumn("click_num", gorm.Expr("click_num + ?", 1)).Error
}

// ============ SKU 相关查询 ============

// 5. 根据商品ID查询所有 SKU
func (d *ProductDao) GetSkusByProductID(productID uint) (skus []*model.ProductSku, err error) {
	err = DB.Model(&model.ProductSku{}).Where("product_id = ?", productID).Find(&skus).Error
	return
}

// 6. 根据 SKU ID 查询单个 SKU
func (d *ProductDao) GetSkuByID(skuID uint) (sku *model.ProductSku, err error) {
	err = DB.Model(&model.ProductSku{}).Where("id = ?", skuID).First(&sku).Error
	return
}

// 7. 查询 SKU 库存（优化版）
func (d *ProductDao) GetSkuStock(skuID uint) (stock int, err error) {
	err = DB.Model(&model.ProductSku{}).Where("id = ?", skuID).
		Select("stock").First(&stock).Error
	return
}

// ============ 库存扣减（乐观锁） ============

// 8. 扣减库存（乐观锁版本）- 关键！
func (d *ProductDao) DecrementStock(tx *gorm.DB, skuID uint, quantity int, version int) (rowsAffected int64, err error) {
	// 正确的乐观锁扣减逻辑
	result := tx.Model(&model.ProductSku{}).
		Where("id = ? AND version = ? AND stock >= ?", skuID, version, quantity).
		// ↑ SKU ID     ↑ 版本号      ↑ 库存充足
		Updates(map[string]interface{}{
			"stock":   gorm.Expr("stock - ?", quantity),
			"version": gorm.Expr("version + ?", 1),
		})

		// 如果 RowsAffected = 0，说明：
		// - version 不匹配（被其他请求抢先更新）
		// - stock 不足
		// 需要在 Service 层重试或返回失败
	return result.RowsAffected, result.Error
}

// 9. 回退库存（取消订单/支付超时）
func (d *ProductDao) IncrementStock(skuID uint, quantity int) error {
	return DB.Model(&model.ProductSku{}).Where("id = ?", skuID).
		Updates(map[string]interface{}{
			"stock":   gorm.Expr("stock + ?", quantity),
			"version": gorm.Expr("version + ?", 1),
		}).Error
}

// ============ 商品管理（CRUD） ============

// internal/dao/product.go

// 10. 创建商品SPU（支持事务）
func (d *ProductDao) CreateProductSPU(tx *gorm.DB, product *model.Product) error {
	return tx.Create(product).Error
	//     ↑ 使用传入的 tx，而不是全局 DB
}

// 11. 创建商品SKU（支持事务）
func (d *ProductDao) CreateProductSKUs(tx *gorm.DB, skus []*model.ProductSku) error {
	return tx.Create(skus).Error
}

// 12. 更新SKU库存（直接设置库存值）
func (d *ProductDao) UpdateSkuStock(skuID uint, stock int) error {
	return DB.Model(&model.ProductSku{}).
		Where("id = ?", skuID).
		Updates(map[string]interface{}{
			"stock":   stock,
			"version": gorm.Expr("version + ?", 1),
		}).Error
}

// 13. 更新商品状态
func (d *ProductDao) UpdateProductOnSale(productID uint, onSale bool) error {
	return DB.Model(&model.Product{}).Where("id=?", productID).Update("on_sale", onSale).Error
}
