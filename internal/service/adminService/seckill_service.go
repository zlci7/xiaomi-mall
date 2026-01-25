package adminService

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	parseTime "xiaomi-mall/pkg/parsetime"
	"xiaomi-mall/pkg/xerr"
)

type SeckillService struct{}

var Seckill = new(SeckillService)

// 创建秒杀商品
func (s *SeckillService) CreateSeckillProduct(req dto.CreateSeckillProductReq) (*vo.CreateSeckillProductResp, error) {

	//1.查询商品是否存在
	product, err := dao.Product.GetProductByID(req.ProductID)
	if err != nil {
		return nil, xerr.NewErrMsg("商品不存在")
	}
	//2.查询商品库存是否存在
	sku, err := dao.Product.GetSkuByID(req.SkuID)
	if err != nil {
		return nil, xerr.NewErrMsg("商品库存不存在")
	}
	if sku.Stock < int(req.SeckillStock) {
		return nil, xerr.NewErrMsg("商品库存不足")
	}

	//3.解析时间
	startTime, err := parseTime.ParseDateTimeStr(req.StartTime)
	endTime, err := parseTime.ParseDateTimeStr(req.EndTime)

	// 4. 创建秒杀商品
	seckill := &model.SeckillProduct{
		ProductID:    product.ID,
		SkuID:        req.SkuID,
		SeckillPrice: req.SeckillPrice,
		SeckillStock: req.SeckillStock,
		TotalStock:   req.SeckillStock,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       0,
		Version:      0,
	}
	err = dao.Seckill.CreateSeckillProduct(seckill)
	if err != nil {
		return nil, xerr.NewErrMsg("创建秒杀商品失败")
	}
	return &vo.CreateSeckillProductResp{
		ID: seckill.ID,
	}, nil
}

// 删除秒杀商品
func (s *SeckillService) DeleteSeckillProduct(req dto.DeleteSeckillProductReq) error {
	return dao.Seckill.DeleteSeckillProduct(req.ID)
}

// 手动开启/结束秒杀
func (s *SeckillService) UpdateSeckillStatus(req dto.UpdateSeckillStatusReq) error {
	return dao.Seckill.UpdateSeckillStatus(req.ID, req.Status)
}
