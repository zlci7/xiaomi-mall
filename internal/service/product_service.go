package service

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
)

type ProductService struct{}

var Product = new(ProductService)

// 创建商品
func (s *ProductService) CreateProduct(req dto.CreateProductReq) (resp *vo.CreateProductResp, err error) {
	return
}
