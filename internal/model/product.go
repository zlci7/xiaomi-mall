package model

import "gorm.io/gorm"

// Category 商品分类
type Category struct {
	gorm.Model
	Name     string `json:"name"`
	ParentID uint   `json:"parent_id"` // 父分类ID，0代表顶级分类
}

// Carousel 轮播图 (首页广告)
// type Carousel struct {
// 	gorm.Model
// 	ImgPath   string `json:"img_path"`
// 	ProductID uint   `json:"product_id"` // 点击跳转到哪个商品
// }

// Product (SPU) 商品主表
type Product struct {
	gorm.Model
	Name          string `gorm:"size:255;index" json:"name"` // 商品名
	CategoryID    uint   `gorm:"not null;index" json:"category_id"`
	Title         string `json:"title"`
	Info          string `gorm:"size:1000" json:"info"` // 详细描述
	ImgPath       string `json:"img_path"`
	Price         int64  `json:"price"`                        // 展示价格，单位：分
	DiscountPrice int64  `json:"discount_price"`               // 折扣价，单位：分
	OnSale        bool   `gorm:"default:false" json:"on_sale"` // 是否上架
	Num           int    `json:"num"`                          // 销量
	ClickNum      int    `json:"click_num"`                    // 点击量
}

// ProductSku (SKU) 商品规格表 —— 库存管理的原子单位
type ProductSku struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index" json:"product_id"`
	Title     string `json:"title"`                       // 规格名，如 "红色+64G"
	Price     int64  `json:"price"`                       // 价格，单位：分
	Stock     int    `gorm:"check:stock>=0" json:"stock"` // 库存，数据库层面约束不能小于0
	Code      string `json:"code"`                        // 商家编码
	Version   int    `gorm:"default:0" json:"version"`    // 乐观锁版本号
}
