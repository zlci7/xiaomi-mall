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

	// ========== 检查事务是否成功 ==========
	if err != nil {
		return nil, err
	}

	// ========== 【事务后】Step 8: 加入延迟队列 ==========
	score := float64(order.ExpireTime.Unix())
	dao.Rdb.ZAdd(ctx, "order:delay:queue", &redis.Z{
		Score:  score,
		Member: orderNum,
	})

	// ========== 【事务后】Step 9: 返回订单信息 ==========
	return &vo.CreateOrderResp{
		OrderNo:     orderNum,
		TotalAmount: totalAmount,
		ExpireTime:  order.ExpireTime,
	}, nil
}

// PayOrder 支付订单（模拟支付）
func (s *OrderService) PayOrder(userID uint, req dto.PayOrderReq) (*vo.PayOrderResp, error) {
	orderNo := req.OrderNo

	// ========== Step 1: 查询订单 ==========
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		return nil, xerr.NewErrMsg("订单不存在")
	}

	// ========== Step 2: 权限校验 ==========
	if order.UserID != userID {
		return nil, xerr.NewErrMsg("订单不属于当前用户")
	}

	// ========== Step 3: 状态校验 ==========
	// 3.1 检查是否已支付
	if order.PayStatus == 1 {
		return nil, xerr.NewErrMsg("订单已支付，请勿重复支付")
	}

	// 3.2 检查订单是否已取消
	if order.OrderStatus == 4 {
		return nil, xerr.NewErrMsg("订单已取消，无法支付")
	}

	// 3.3 检查订单是否已过期
	if time.Now().After(order.ExpireTime) {
		return nil, xerr.NewErrMsg("订单已过期")
	}

	// ========== Step 4: 模拟支付成功 ==========
	// 真实场景需要调用支付宝/微信 SDK，获取支付二维码
	// 这里模拟直接支付成功
	tradeNo := generateMockTradeNo(req.PayType) // 模拟交易流水号

	// ========== Step 5: 更新订单状态（事务 + 乐观锁） ==========
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 使用 DAO 层封装的支付方法（包含乐观锁）
		rowsAffected, err := dao.Order.PayOrder(
			tx,
			orderNo,
			req.PayType,
			tradeNo,
			order.Version,
		)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return xerr.NewErrMsg("订单状态已变更，请刷新后重试")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// ========== Step 6: 从延迟队列移除 ==========
	// 已支付的订单不需要自动关闭
	dao.Rdb.ZRem(ctx, "order:delay:queue", orderNo)

	// ========== Step 7: 返回支付结果 ==========
	return &vo.PayOrderResp{
		OrderNo:   orderNo,
		PayStatus: 1, // 已支付
		TradeNo:   tradeNo,
	}, nil
}

// generateMockTradeNo 生成模拟的支付流水号
func generateMockTradeNo(payType int) string {
	// 真实场景：返回支付宝/微信的交易流水号
	// 模拟场景：使用雪花算法生成唯一流水号
	return idgen.GenStringID()
}

// CloseOrder 自动关闭订单（超时未支付）
func (s *OrderService) CloseOrder(orderNo string) error {
	// 1️⃣ 查询订单
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		// 订单不存在（可能是之前创建失败的脏数据），直接跳过
		if err == gorm.ErrRecordNotFound {
			return nil
		}
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

// 取消订单
func (s *OrderService) CancelOrder(userID uint, req dto.CancelOrderReq) error {
	orderNo := req.OrderNo
	// 1️⃣ 查询订单
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return xerr.NewErrMsg("订单不属于当前用户")
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

// 订单详情查询
func (s *OrderService) GetOrderDetail(req dto.OrderDetailReq) (*vo.OrderDetailResp, error) {
	orderNo := req.OrderNo

	// ========== Step 1: 查询订单主表 ==========
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		return nil, xerr.NewErrMsg("订单不存在")
	}

	// ========== Step 2: 查询订单商品列表 ==========
	orderItems, err := dao.Order.GetOrderItems(orderNo)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// ========== Step 3: 解析地址快照 ==========
	var addressSnapshot vo.AddressSnapshotVO
	if err := json.Unmarshal([]byte(order.AddressSnapshot), &addressSnapshot); err != nil {
		// 如果解析失败，返回空地址（兼容旧数据）
		addressSnapshot = vo.AddressSnapshotVO{}
	}

	// ========== Step 4: 组装订单商品列表 ==========
	items := make([]vo.OrderDetailItemVO, 0, len(orderItems))
	for _, item := range orderItems {
		items = append(items, vo.OrderDetailItemVO{
			ProductID:    item.ProductID,
			ProductSkuID: item.ProductSkuID,
			Title:        item.Title,
			ImgPath:      item.ImgPath,
			Price:        item.Price,
			Num:          item.Num,
			Subtotal:     item.Price * int64(item.Num), // 小计 = 单价 * 数量
		})
	}

	// ========== Step 5: 组装响应 ==========
	// 注意：时间字段已是指针类型，无需转换
	resp := &vo.OrderDetailResp{
		// 订单基本信息
		OrderNo:     order.OrderNum,
		TotalAmount: order.AllPrice,
		OrderStatus: order.OrderStatus,
		PayStatus:   order.PayStatus,
		PayType:     order.PayType,
		Type:        order.Type,

		// 时间信息（已是指针类型，直接赋值）
		CreatedAt:  order.CreatedAt,
		PayTime:    order.PayTime,
		ShipTime:   order.ShipTime,
		FinishTime: order.FinishTime,
		CancelTime: order.CancelTime,
		ExpireTime: order.ExpireTime,

		// 收货地址
		Address: addressSnapshot,

		// 商品列表
		Items: items,

		// 物流信息
		TrackingNumber: order.TrackingNumber,

		// 备注
		Remark:      order.Remark,
		AdminRemark: order.AdminRemark,
	}

	return resp, nil
}

// 订单列表查询
func (s *OrderService) GetOrderList(userID uint, req dto.OrderListReq) (*vo.OrderListResp, error) {
	// ========== Step 1: 设置默认分页参数 ==========
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// ========== Step 2: 查询订单列表 ==========
	orders, total, err := dao.Order.GetUserOrders(userID, page, pageSize, req.OrderStatus)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// ========== Step 3: 组装订单列表 VO ==========
	list := make([]vo.OrderItemVO, 0, len(orders))
	for _, order := range orders {
		// 查询该订单的商品列表（用于获取第一个商品和总数量）
		orderItems, err := dao.Order.GetOrderItems(order.OrderNum)
		if err != nil {
			continue // 跳过异常订单
		}

		// 计算商品总数量
		productCount := 0
		for _, item := range orderItems {
			productCount += item.Num
		}

		// 获取第一个商品（用于列表展示）
		var firstProduct vo.ProductSnapshotVO
		if len(orderItems) > 0 {
			firstItem := orderItems[0]
			firstProduct = vo.ProductSnapshotVO{
				Title:   firstItem.Title,
				ImgPath: firstItem.ImgPath,
				Price:   firstItem.Price,
				Num:     firstItem.Num,
			}
		}

		list = append(list, vo.OrderItemVO{
			OrderNo:      order.OrderNum,
			TotalAmount:  order.AllPrice,
			OrderStatus:  order.OrderStatus,
			PayStatus:    order.PayStatus,
			ProductCount: productCount,
			FirstProduct: firstProduct,
			CreatedAt:    order.CreatedAt,
			ExpireTime:   order.ExpireTime,
		})
	}

	// ========== Step 4: 返回分页结果 ==========
	return &vo.OrderListResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 确认收货
func (s *OrderService) ConfirmOrder(userID uint, req dto.ConfirmOrderReq) error {
	orderNo := req.OrderNo

	// ========== Step 1: 查询订单 ==========
	order, err := dao.Order.GetOrderByOrderNum(orderNo)
	if err != nil {
		return xerr.NewErrMsg("订单不存在")
	}

	// ========== Step 2: 权限校验 ==========
	if order.UserID != userID {
		return xerr.NewErrMsg("订单不属于当前用户")
	}

	// ========== Step 3: 状态校验 ==========
	// 3.1 检查支付状态
	if order.PayStatus != 1 {
		return xerr.NewErrMsg("订单未支付，无法确认收货")
	}

	// 3.2 检查订单状态（只有已发货的订单才能确认收货）
	if order.OrderStatus != 2 {
		if order.OrderStatus == 0 {
			return xerr.NewErrMsg("订单未支付")
		} else if order.OrderStatus == 1 {
			return xerr.NewErrMsg("订单未发货，无法确认收货")
		} else if order.OrderStatus == 3 {
			return xerr.NewErrMsg("订单已完成")
		} else if order.OrderStatus == 4 {
			return xerr.NewErrMsg("订单已取消")
		}
		return xerr.NewErrMsg("订单状态异常")
	}

	// ========== Step 4: 更新订单状态（事务 + 乐观锁） ==========
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		rowsAffected, err := dao.Order.ConfirmOrder(tx, orderNo, order.Version)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return xerr.NewErrMsg("订单状态已变更，请刷新后重试")
		}
		return nil
	})

	return err
}
