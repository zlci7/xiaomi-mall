// internal/service/userService/seckill_service.go

package userService

import (
	"context"
	"encoding/json"
	"time"

	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/pkg/xerr"
)

type SeckillService struct{}

var Seckill = new(SeckillService)

// Redis 缓存数据结构（与预热时一致）
type SeckillProductCache struct {
	SeckillID     uint   `json:"seckill_id"`
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductTitle  string `json:"title"`
	ProductInfo   string `json:"info"`
	ProductImg    string `json:"img"`
	CategoryID    uint   `json:"category_id"`
	SkuID         uint   `json:"sku_id"`
	SkuTitle      string `json:"sku_title"`
	SkuCode       string `json:"sku_code"`
	OriginalPrice int64  `json:"original_price"`
	SeckillPrice  uint   `json:"seckill_price"`
	SeckillStock  uint   `json:"seckill_stock"` // 总库存
	TotalStock    uint   `json:"total_stock"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
}

// ==================== 列表查询 ====================

func (s *SeckillService) GetSeckillProductList(req dto.UserSeckillListReq) (*vo.UserSeckillListResp, error) {
	ctx := context.Background()

	// 1. 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 2. 从 Redis 获取未结束的秒杀ID（分页）
	seckillIDs, total, err := dao.Seckill.GetActiveSeckillIDs(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewErrMsg("查询秒杀活动失败")
	}

	if len(seckillIDs) == 0 {
		return &vo.UserSeckillListResp{
			List:  []vo.UserSeckillListItemVO{},
			Total: 0,
		}, nil
	}

	// 3. 批量获取商品详情
	productDataMap, err := dao.Seckill.BatchGetSeckillProductCache(ctx, seckillIDs)
	if err != nil {
		return nil, xerr.NewErrMsg("获取商品详情失败")
	}

	// 4. 批量获取实时库存
	stockMap, err := dao.Seckill.BatchGetSeckillStocks(ctx, seckillIDs)
	if err != nil {
		return nil, xerr.NewErrMsg("获取库存失败")
	}

	// 5. 组装 VO
	now := time.Now()
	nowUnix := now.Unix()
	list := make([]vo.UserSeckillListItemVO, 0, len(seckillIDs))

	for _, id := range seckillIDs {
		// 解析商品详情
		productData, ok := productDataMap[id]
		if !ok {
			continue
		}

		var cache SeckillProductCache
		if err := json.Unmarshal(productData, &cache); err != nil {
			continue
		}

		// 获取实时库存
		stock := stockMap[id]

		// 计算已售数量
		soldNum := int(cache.TotalStock) - stock
		if soldNum < 0 {
			soldNum = 0
		}

		// 判断状态
		var status string
		var canBuy bool
		isSoldOut := stock <= 0

		if nowUnix < cache.StartTime {
			status = "未开始"
			canBuy = false
		} else if nowUnix >= cache.StartTime && nowUnix < cache.EndTime {
			status = "进行中"
			canBuy = !isSoldOut
		} else {
			status = "已结束"
			canBuy = false
		}

		// 组装 VO
		voItem := vo.UserSeckillListItemVO{
			ID:            cache.SeckillID,
			ProductID:     cache.ProductID,
			ProductName:   cache.ProductName,
			ImgPath:       cache.ProductImg,
			OriginalPrice: uint(cache.OriginalPrice),
			SeckillPrice:  cache.SeckillPrice,
			SeckillStock:  uint(stock),
			SoldNum:       uint(soldNum),
			StartTime:     time.Unix(cache.StartTime, 0),
			EndTime:       time.Unix(cache.EndTime, 0),
			Status:        status,
			IsSoldOut:     isSoldOut,
			CanBuy:        canBuy,
		}

		list = append(list, voItem)
	}

	return &vo.UserSeckillListResp{
		List:  list,
		Total: total,
	}, nil
}

// ==================== 详情查询 ====================

func (s *SeckillService) GetSeckillProductDetail(userID uint, req dto.UserSeckillDetailReq) (*vo.UserSeckillDetailVO, error) {
	ctx := context.Background()

	// 1. 从 Redis 获取商品详情
	productData, err := dao.Seckill.GetSeckillProductCacheByID(ctx, req.ID)
	if err != nil {
		return nil, xerr.NewErrMsg("秒杀商品不存在或已结束")
	}

	var cache SeckillProductCache
	if err := json.Unmarshal(productData, &cache); err != nil {
		return nil, xerr.NewErrMsg("数据解析失败")
	}

	// 2. 获取实时库存
	stock, err := dao.Seckill.GetSeckillStockByID(ctx, req.ID)
	if err != nil {
		return nil, xerr.NewErrMsg("获取库存失败")
	}

	// 3. 检查用户是否已购买（如果传了 userID）
	hasPurchased := false
	if userID > 0 {
		hasPurchased, _ = dao.Seckill.CheckUserPurchased(ctx, req.ID, userID)
	}

	// 4. 计算已售数量
	soldNum := int(cache.TotalStock) - stock
	if soldNum < 0 {
		soldNum = 0
	}

	// 5. 判断状态
	now := time.Now()
	nowUnix := now.Unix()

	var status string
	var canBuy bool
	isSoldOut := stock <= 0

	if nowUnix < cache.StartTime {
		status = "未开始"
		canBuy = false
	} else if nowUnix >= cache.StartTime && nowUnix < cache.EndTime {
		status = "进行中"
		canBuy = !isSoldOut && !hasPurchased
	} else {
		status = "已结束"
		canBuy = false
	}

	// 6. 组装详情 VO
	detailVO := &vo.UserSeckillDetailVO{
		ID:            cache.SeckillID,
		ProductID:     cache.ProductID,
		ProductName:   cache.ProductName,
		ProductTitle:  cache.ProductTitle,
		ProductInfo:   cache.ProductInfo,
		ImgPath:       cache.ProductImg,
		CategoryID:    cache.CategoryID,
		SkuID:         cache.SkuID,
		SkuTitle:      cache.SkuTitle,
		SkuCode:       cache.SkuCode,
		OriginalPrice: uint(cache.OriginalPrice),
		SeckillPrice:  cache.SeckillPrice,
		SeckillStock:  uint(stock),
		TotalStock:    cache.TotalStock,
		SoldNum:       uint(soldNum),
		StartTime:     time.Unix(cache.StartTime, 0),
		EndTime:       time.Unix(cache.EndTime, 0),
		Status:        status,
		IsSoldOut:     isSoldOut,
		CanBuy:        canBuy,
		HasPurchased:  hasPurchased,
	}

	return detailVO, nil
}
