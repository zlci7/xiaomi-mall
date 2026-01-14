package vo

// 商品列表项（简化版）
type ProductItemVO struct {
	ProductID     uint   `json:"product_id"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	ImgPath       string `json:"img_path"`
	Price         int64  `json:"price"`
	DiscountPrice int64  `json:"discount_price"`
	Num           int    `json:"num"`       // 销量
	ClickNum      int    `json:"click_num"` // 点击量
	OnSale        bool   `json:"on_sale"`
}

// SKU VO
type SkuVO struct {
	SkuID uint   `json:"sku_id"`
	Title string `json:"title"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
	Code  string `json:"code"`
}

// 商品分类VO
type CategoryVO struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
}

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
// type UpdateProductStockResp struct {
// 	ProductID   uint   `json:"product_id"`
// 	ProductName string `json:"product_name"`
// 	Stock       int    `json:"stock"`
// }

// 更新商品上架状态请求
// type UpdateProductOnSaleResp struct {
// 	ProductID uint `json:"product_id"`
// 	OnSale    bool `json:"on_sale"`
// }

// ============ 商品查询 VO ============

// 商品列表响应
type ProductListResp struct {
	List     []ProductItemVO `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// 商品详情响应（完整版）
type ProductDetailResp struct {
	ProductID     uint    `json:"product_id"`
	Name          string  `json:"name"`
	CategoryID    uint    `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	Title         string  `json:"title"`
	Info          string  `json:"info"`
	ImgPath       string  `json:"img_path"`
	Price         int64   `json:"price"`
	DiscountPrice int64   `json:"discount_price"`
	Num           int     `json:"num"`
	ClickNum      int     `json:"click_num"`
	OnSale        bool    `json:"on_sale"`
	SKUs          []SkuVO `json:"skus"` // ⬅️ 包含 SKU 列表
}

// SKU详情响应
type SkuDetailResp struct {
	SkuID uint   `json:"sku_id"`
	Title string `json:"title"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
	Code  string `json:"code"`
}

// 商品分类列表响应
type CategoryListResp struct {
	List []CategoryVO `json:"list"`
}
