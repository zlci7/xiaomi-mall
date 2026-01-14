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
	Page       int    `form:"page" binding:"omitempty,min=1"`                                      // 页码，默认1
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`                         // 每页数量，默认10，最大100
	CategoryID uint   `form:"category_id"`                                                         // 分类ID，可选
	Keyword    string `form:"keyword" binding:"omitempty,max=100"`                                 // 搜索关键词，可选，最大100字符
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=price num click_num created_at id"` // 排序字段：price(价格), num(销量), click_num(点击量), created_at(创建时间)，默认id
	Order      string `form:"order" binding:"omitempty,oneof=asc desc"`                            // 排序方向：asc(升序), desc(降序)，默认desc
	OnSale     *bool  `form:"on_sale"`                                                             // 是否上架，可选（nil表示不筛选）
}

// 商品详情请求 - GET /products/:product_id
type ProductDetailReq struct {
	ProductID uint `uri:"product_id" binding:"required,min=1"`
}

// SKU 详情请求 - GET /products/skus/:sku_id
type SkuDetailReq struct {
	SkuID uint `uri:"sku_id" binding:"required,min=1"`
}
