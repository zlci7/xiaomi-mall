// internal/service/userService/seckill_service.go

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
	"xiaomi-mall/internal/pkg/types"
	"xiaomi-mall/pkg/idgen"
	"xiaomi-mall/pkg/xerr"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SeckillService struct{}

var Seckill = new(SeckillService)

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

		var cache types.SeckillProductCache
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

	var cache types.SeckillProductCache
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

// ==================== 秒杀下单 ====================

func (s *SeckillService) CreateSeckillOrder(userID uint, req dto.CreateSeckillOrderReq) (*vo.CreateSeckillOrderResp, error) {
	ctx := context.Background()
	//1.1 检查商品是否在redis
	cacheData, err := dao.Seckill.GetSeckillProductCacheByID(ctx, req.SeckillProductID)
	if err != nil {
		return nil, xerr.NewErrMsg("秒杀商品不存在")
	}

	//1.2 检查当前是否在秒杀活动时间内
	var productCache types.SeckillProductCache
	json.Unmarshal(cacheData, &productCache)
	now := time.Now().Unix()
	if now < productCache.StartTime || now >= productCache.EndTime {
		return nil, xerr.NewErrMsg("秒杀活动未开始或已结束")
	}

	//1.3 检查用户是否已购买（Redis，快速检查）
	hasBought, err := dao.Seckill.CheckUserPurchased(ctx, req.SeckillProductID, userID)
	if err != nil {
		return nil, xerr.NewErrMsg("系统错误，请稍后重试")
	}
	if hasBought {
		return nil, xerr.NewErrMsg("请勿重复下单")
	}

	//2.1 执行lua脚本
	//生成订单号
	orderNum := idgen.GenStringID()
	result, err := dao.Seckill.CreateSeckillOrderAtomic(ctx, userID, req.SeckillProductID, orderNum)
	if err != nil {
		return nil, xerr.NewErrMsg("系统错误，请稍后重试")
	}
	if result == -1 {
		return nil, xerr.NewErrMsg("请勿重复下单")
	}
	if result == -2 {
		return nil, xerr.NewErrMsg("库存不足")
	}
	if result == -3 {
		return nil, xerr.NewErrMsg("秒杀活动未开始或已结束")
	}

	//2.2 将订单加入延迟队列（30分钟后超时）
	expireTime := time.Now().Add(30 * time.Minute)
	dao.Rdb.ZAdd(ctx, "order:delay:queue", &redis.Z{
		Score:  float64(expireTime.Unix()),
		Member: orderNum,
	})

	//2.3 返回秒杀成功（异步写入MySQL已在Lua脚本中投递到队列）
	return &vo.CreateSeckillOrderResp{
		OrderNum:    orderNum,
		TotalAmount: int64(productCache.SeckillPrice),
		ExpireTime:  expireTime,
		PayUrl:      "",
	}, nil
}

// ==================== 秒杀订单关闭（超时取消）====================

// CloseSeckillOrder 关闭秒杀订单（回滚 Redis 库存和用户标记）
func (s *SeckillService) CloseSeckillOrder(orderNum string) error {
	ctx := context.Background()

	// 1. 查询主订单信息
	var order model.Order
	if err := dao.DB.Where("order_num = ?", orderNum).First(&order).Error; err != nil {
		return err
	}

	// 2. 检查订单状态（只关闭未支付的订单）
	if order.PayStatus != 0 {
		return nil // 已支付，跳过
	}

	// 3. 查询秒杀订单信息
	var seckillOrder model.SeckillOrder
	if err := dao.DB.Where("order_num = ?", orderNum).First(&seckillOrder).Error; err != nil {
		return err
	}

	// 4. 事务：更新数据库订单状态
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		// 4.1 更新主订单状态
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"order_status": 4, // 4=已取消
			"cancel_time":  time.Now(),
		}).Error; err != nil {
			return err
		}

		// 4.2 更新秒杀订单状态
		if err := tx.Model(&seckillOrder).Update("status", 2).Error; err != nil { // 2=已取消
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 5. 回滚 Redis 库存（关键！）
	stockKey := fmt.Sprintf("seckill:stock:%d", seckillOrder.SeckillProductID)
	dao.Rdb.Incr(ctx, stockKey)

	// 6. 删除用户购买标记（关键！）
	userKey := fmt.Sprintf("seckill:user:%d:%d", seckillOrder.SeckillProductID, order.UserID)
	dao.Rdb.Del(ctx, userKey)

	return nil
}
