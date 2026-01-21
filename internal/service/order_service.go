package service

import (
	"encoding/json"
	"time"
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	"xiaomi-mall/pkg/idgen"
	"xiaomi-mall/pkg/xerr"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderService struct{}

var Order = new(OrderService)

// 创建订单
func (OrderService) CreateOrder(userID uint, req dto.CreateOrderReq) (resp *vo.CreateOrderResp, err error) {

	// ========== 【事务外】Step 1: 参数校验 ==========
	if len(req.Items) == 0 {
		return nil, xerr.NewErrCode(xerr.REUQEST_PARAM_ERROR)
	}
	// ========== 【事务外】Step 2: 批量查询 SKU 信息 ==========
	skuIDs := make([]uint, 0, len(req.Items))
	for _, item := range req.Items {
		skuIDs = append(skuIDs, item.SkuID)
	}
	skus, err := dao.Product.GetSkusByIDs(skuIDs)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}
	//用map存方便查询，key是skuID，value是sku指针，为什么存指针是为了方便修改
	skuMap := make(map[uint]*model.ProductSku)
	for _, sku := range skus {
		skuMap[sku.ID] = sku
	}
	for _, item := range req.Items {
		sku, exists := skuMap[item.SkuID]
		if !exists {
			return nil, xerr.NewErrMsg("商品不存在")
		}
		if sku.Stock < item.Num {
			return nil, xerr.NewErrMsg("库存不足")
		}
	}

	// ========== 【事务外】Step 3: 查询用户地址 ==========
	address, err := dao.Address.GetByID(req.AddressID)
	if err != nil {
		return nil, xerr.NewErrMsg("地址不存在")
	}
	if address.UserID != userID {
		return nil, xerr.NewErrMsg("地址不属于当前用户")
	}

	// ========== 【事务外】Step 4: 计算订单总金额 ==========
	totalAmount := int64(0)
	for _, sku := range req.Items {
		totalAmount += skuMap[sku.SkuID].Price * int64(sku.Num)
	}
	// ========== 【事务外】Step 5: 生成订单号（雪花算法）==========
	orderNum := idgen.GenStringID()

	// ========== 【事务外】Step 6: 准备快照数据 ==========
	// 6.1 地址快照
	addressJSON, _ := json.Marshal(map[string]string{
		"name":    address.Name,
		"phone":   address.Phone,
		"address": address.Address,
	})
	// 6.2 订单主表数据
	order := &model.Order{
		UserID:          userID,
		OrderNum:        orderNum,
		AllPrice:        totalAmount,
		PayStatus:       0,
		OrderStatus:     0,
		Type:            1, // 普通订单
		AddressSnapshot: string(addressJSON),
		ExpireTime:      time.Now().Add(30 * time.Minute),
		Remark:          req.Remark,
		Version:         0,
	}

	// 6.3 订单详情数据（商品快照）
	orderItems := make([]*model.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		sku := skuMap[item.SkuID]
		orderItems = append(orderItems, &model.OrderItem{
			OrderNum:     orderNum,
			ProductID:    sku.ProductID,
			ProductSkuID: sku.ID,
			Num:          item.Num,
			Price:        sku.Price,
			Title:        sku.Title,
			ImgPath:      sku.ImgPath,
		})
	}

	// ========== 【事务内】Step 7: 开启数据库事务 ==========
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 7.1 乐观锁扣减库存
		for _, item := range req.Items {
			sku := skuMap[item.SkuID]
			rowsAffected, err := dao.Product.DecrementStock(tx, item.SkuID, item.Num, sku.Version)
			if err != nil {
				return err // 库存不足或版本冲突，回滚
			}
			if rowsAffected == 0 {
				return xerr.NewErrMsg("库存不足")
			}
		}

		// 7.2 创建订单
		if err := dao.Order.CreateOrder(tx, order, orderItems); err != nil {
			return err
		}

		// 7.3 清空购物车（可选）
		// if req.FromCart {
		//     tx.Where("user_id = ? AND sku_id IN ?", userID, skuIDs).
		//       Delete(&model.Cart{})
		// }

		return nil
	})

	// ========== 【事务后】Step 8: 加入延迟队列 ==========
	score := float64(order.ExpireTime.Unix())
	dao.Rdb.ZAdd(ctx, "order:delay:queue", &redis.Z{
		Score:  score,
		Member: orderNum,
	})

	// ========== 【事务后】Step 9: 返回订单信息 ==========\
	return &vo.CreateOrderResp{
		OrderNo:     orderNum,
		TotalAmount: totalAmount,
		ExpireTime:  order.ExpireTime,
	}, nil
}

// CloseOrder 自动关闭订单（超时未支付）
func (s *OrderService) CloseOrder(orderNo string) error {
	// 1️⃣ 查询订单
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		return err
	}

	// 2️⃣ 只关闭未支付的订单
	if order.PayStatus != 0 {
		return nil // 已支付，跳过
	}
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		// 3️⃣ 更新订单状态（乐观锁）
		rowsAffected, err := dao.Order.UpdateOrderStatus(
			tx,
			orderNo,
			4, // OrderStatus = 4 (已取消)
			order.Version,
		)

		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return xerr.NewErrMsg("订单状态已变更")
		}

		// 4️⃣ 回滚库存
		items, err := dao.Order.GetOrderItems(orderNo)
		if err != nil {
			return err
		}

		for _, item := range items {
			tx.Model(&model.ProductSku{}).
				Where("id = ?", item.ProductSkuID).
				Updates(map[string]interface{}{
					"stock":   gorm.Expr("stock + ?", item.Num),
					"version": gorm.Expr("version + 1"),
				})
		}

		return nil
	})
}
