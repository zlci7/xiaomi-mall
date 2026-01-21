package dto

// ========== 创建/编辑地址 ==========
type SaveAddressReq struct {
	ID        uint   `json:"id"` // ID > 0 表示编辑，= 0 表示新增
	Name      string `json:"name" binding:"required,max=50"`
	Phone     string `json:"phone" binding:"required,len=11,numeric"`
	Address   string `json:"address" binding:"required,max=200"`
	IsDefault bool   `json:"is_default"` // 是否默认地址
}

// ========== 删除地址 ==========
type DeleteAddressReq struct {
	ID uint `json:"id" binding:"required,min=1"`
}

// ========== 地址列表 ==========
// 无请求参数，直接查询当前用户的所有地址
