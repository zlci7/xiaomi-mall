package vo

// ========== 地址列表响应 ==========
type AddressListResp struct {
	List []AddressVO `json:"list"`
}

// 地址 VO
type AddressVO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
}
