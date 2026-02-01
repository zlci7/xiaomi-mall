package adminService

import (
	"context"
	"encoding/json"
	"time"
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

// 预热秒杀商品到 Redis
func (s *SeckillService) PreheatSeckillProduct(req dto.PreheatSeckillProductReq) error {
	ctx := context.Background() // ← 添加这行
	//1. 【数据库】查询秒杀商品详情
	//a.查询秒杀商品具体信息
	seckillProduct, err := dao.Seckill.GetSeckillProductByID(req.ID)
	if err != nil {
		return xerr.NewErrMsg("秒杀商品不存在")
	}

	//b.验证状态（未开始）
	if seckillProduct.Status != 0 {
		return xerr.NewErrMsg("无法预热，秒杀活动已开始或结束")
	}

	//c.查询products和sku库存信息
	product, err := dao.Product.GetProductByID(seckillProduct.ProductID)
	if err != nil {
		return xerr.NewErrMsg("商品spu不存在")
	}
	sku, err := dao.Product.GetSkuByID(seckillProduct.SkuID)
	if err != nil {
		return xerr.NewErrMsg("商品sku不存在")
	}
	if sku.Stock < int(seckillProduct.SeckillStock) {
		return xerr.NewErrMsg("商品库存不足")
	}

	//2. 【Redis String】设置库存
	//2. 设置库存（调用 DAO 方法）
	ttl := time.Until(seckillProduct.EndTime) + 24*time.Hour
	ttlSeconds := int64(ttl.Seconds()) // ← 转换为秒
	//3. 【Redis Hash】缓存商品详情
	cacheData := map[string]interface{}{ // ← 改为包含秒杀信息
		"seckill_id":     seckillProduct.ID,
		"product_id":     product.ID,
		"product_name":   product.Name,
		"product_img":    product.ImgPath,
		"title":          product.Title,
		"info":           product.Info,
		"sku_id":         sku.ID,
		"sku_title":      sku.Title,
		"original_price": sku.Price,
		"seckill_price":  seckillProduct.SeckillPrice,
		"seckill_stock":  seckillProduct.SeckillStock,
		"start_time":     seckillProduct.StartTime.Unix(),
		"end_time":       seckillProduct.EndTime.Unix(),
	}
	data, err := json.Marshal(cacheData) // ← 处理错误
	if err != nil {
		return xerr.NewErrMsg("序列化失败")
	}
	err = dao.Seckill.PreheatSeckillAtomic(ctx, seckillProduct, data, int64(ttlSeconds))
	if err != nil {
		return xerr.NewErrMsg("预热失败: " + err.Error())
	}

	//4. 【布隆过滤器】添加商品 ID（可选）

	//5. 【数据库】更新预热状态（可选）
	return nil
}
