package vo

// ============ 商品管理 VO ============
// 创建商品请求
type CreateProductResp struct {
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	CategoryID    uint   `json:"category_id"`
	CategoryName  string `json:"category_name"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         int64  `json:"price"`
	DiscountPrice int64  `json:"discount_price"`
}

// 更新商品库存请求
type UpdateProductStockResp struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

// 更新商品上架状态请求
type UpdateProductOnSaleResp struct {
	ProductID uint `json:"product_id"`
	OnSale    bool `json:"on_sale"`
}

// ============ 商品查询 VO ============
