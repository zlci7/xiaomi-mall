package dto

type SKUItem struct {
	Title string `json:"title" binding:"required"`
	Price int64  `json:"price" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
	Code  string `json:"code"`
}

// ============ 商品管理 DTO ============
// 创建商品请求

type CreateProductReq struct {
	Name          string    `json:"name" binding:"required"`
	CategoryID    uint      `json:"category_id" binding:"required"`
	Title         string    `json:"title"`
	Info          string    `json:"info"`
	ImgPath       string    `json:"img_path"`
	Price         int64     `json:"price"`
	DiscountPrice int64     `json:"discount_price"`
	SKUs          []SKUItem `json:"skus" binding:"required,min=1"`
}

// 更新商品库存请求
type UpdateProductStockReq struct {
	ProductID    uint `json:"product_id"`
	Stock        int  `json:"stock"`
	ProductSKUID uint `json:"product_sku_id"`
}

// 更新商品上架状态请求
type UpdateProductOnSaleReq struct {
	ProductID uint `json:"product_id"`
	OnSale    bool `json:"on_sale"`
}

// ============ 商品查询 DTO ============

type ProductListReq struct {
}

type ProductDetailReq struct {
}
