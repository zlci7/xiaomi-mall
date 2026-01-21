package service

import (
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	"xiaomi-mall/pkg/xerr"
)

type AddressService struct{}

var Address = new(AddressService)

// ========== 获取用户地址列表 ==========
func (s *AddressService) GetAddressList(userID uint) (*vo.AddressListResp, error) {
	// 查询用户的所有地址
	addresses, err := dao.Address.GetUserAddresses(userID)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	// 转换为 VO
	addressVOs := make([]vo.AddressVO, 0, len(addresses))
	for _, addr := range addresses {
		addressVOs = append(addressVOs, vo.AddressVO{
			ID:        addr.ID,
			Name:      addr.Name,
			Phone:     addr.Phone,
			Address:   addr.Address,
			IsDefault: addr.IsDefault,
		})
	}

	return &vo.AddressListResp{
		List: addressVOs,
	}, nil
}

// ========== 保存地址（创建/编辑）==========
func (s *AddressService) SaveAddress(userID uint, req dto.SaveAddressReq) error {
	// ========== 新增地址 ==========
	if req.ID == 0 {
		// 1. 检查地址数量（限制最多 10 个）
		count, _ := dao.Address.CountUserAddresses(userID)
		if count >= 10 {
			return xerr.NewErrMsg("最多只能添加10个地址")
		}

		// 2. 如果是默认地址，先取消其他默认地址
		if req.IsDefault {
			dao.Address.CancelDefaultAddress(userID)
		}

		// 3. 如果是第一个地址，自动设为默认
		if count == 0 {
			req.IsDefault = true
		}

		// 4. 创建地址
		address := &model.Address{
			UserID:    userID,
			Name:      req.Name,
			Phone:     req.Phone,
			Address:   req.Address,
			IsDefault: req.IsDefault,
		}

		if err := dao.Address.CreateAddress(address); err != nil {
			return xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
		}

		return nil
	}

	// ========== 编辑地址 ==========
	// 1. 查询地址是否存在
	address, err := dao.Address.GetByID(req.ID)
	if err != nil {
		return xerr.NewErrMsg("地址不存在")
	}

	// 2. 校验地址是否属于该用户
	if address.UserID != userID {
		return xerr.NewErrMsg("无权限操作该地址")
	}

	// 3. 如果设为默认地址，先取消其他默认地址
	if req.IsDefault && !address.IsDefault {
		dao.Address.CancelDefaultAddress(userID)
	}

	// 4. 更新地址
	address.Name = req.Name
	address.Phone = req.Phone
	address.Address = req.Address
	address.IsDefault = req.IsDefault

	if err := dao.Address.UpdateAddress(address); err != nil {
		return xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	return nil
}

// ========== 删除地址 ==========
func (s *AddressService) DeleteAddress(userID uint, req dto.DeleteAddressReq) error {
	// 1. 查询地址是否存在
	address, err := dao.Address.GetByID(req.ID)
	if err != nil {
		return xerr.NewErrMsg("地址不存在")
	}

	// 2. 校验地址是否属于该用户
	if address.UserID != userID {
		return xerr.NewErrMsg("无权限操作该地址")
	}

	// 3. 不允许删除默认地址（需要先设置其他地址为默认）
	if address.IsDefault {
		return xerr.NewErrMsg("请先设置其他地址为默认地址")
	}

	// 4. 删除地址
	if err := dao.Address.DeleteAddress(req.ID); err != nil {
		return xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	return nil
}
