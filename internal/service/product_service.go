package service

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	"xiaomi-mall/pkg/xerr"

	"gorm.io/gorm"
)

type ProductService struct{}

var Product = new(ProductService)

// CreateProduct 添加商品（SPU + SKU）- 使用事务
func (s *ProductService) CreateProduct(req dto.CreateProductReq) (*vo.CreateProductResp, error) {
	// 1️⃣ 准备 SPU 数据
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       req.ImgPath,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        false, // 默认不上架
		Num:           0,
		ClickNum:      0,
	}

	// 2️⃣ 开启事务
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 2.1 创建 SPU
		if err := dao.Product.CreateProductSPU(tx, product); err != nil {
			return err // 返回错误会自动回滚
		}

		// 2.2 准备 SKU 数据（需要关联 ProductID）
		skus := make([]*model.ProductSku, 0, len(req.SKUs))
		for _, skuReq := range req.SKUs {
			sku := &model.ProductSku{
				ProductID: product.ID, // ← 关联刚创建的 SPU ID
				Title:     skuReq.Title,
				Price:     skuReq.Price,
				Stock:     skuReq.Stock,
				Code:      skuReq.Code,
				Version:   0, // 初始版本号
			}
			skus = append(skus, sku)
		}

		// 2.3 批量创建 SKU
		if err := dao.Product.CreateProductSKUs(tx, skus); err != nil {
			return err
		}

		// 2.4 如果需要初始化 Redis 库存，也在这里做
		// for _, sku := range skus {
		//     redisKey := fmt.Sprintf("sku:stock:%d", sku.ID)
		//     dao.Redis.Set(redisKey, sku.Stock, 0)
		// }

		return nil // 返回 nil 表示成功，自动提交事务
	})

	// 3️⃣ 事务失败处理
	if err != nil {
		return nil, xerr.NewErrCode(xerr.PRODUCT_CREATE_ERROR)
	}

	// 4️⃣ 构造响应 VO
	resp := &vo.CreateProductResp{
		ProductID:     product.ID,
		ProductName:   product.Name,
		CategoryID:    product.CategoryID,
		Title:         product.Title,
		Info:          product.Info,
		ImgPath:       product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
	}

	return resp, nil
}
