package dao

import "xiaomi-mall/internal/model"

var Category = new(CategoryDao)

type CategoryDao struct{}

// 1. 获取所有分类（简化版，只有一级分类）
func (d *CategoryDao) GetAllCategories() (categories []*model.Category, err error) {
	// 注意：Category 模型中没有 parent_id 字段，所以是扁平结构
	err = DB.Model(&model.Category{}).
		Order("id ASC").
		Find(&categories).Error
	return
}

// 2. 根据ID查询单个分类
func (d *CategoryDao) GetCategoryByID(categoryID uint) (category *model.Category, err error) {
	err = DB.Model(&model.Category{}).Where("id = ?", categoryID).First(&category).Error
	return
}

// 3. 根据分类查询商品数量（用于展示）
func (d *CategoryDao) CountProductsByCategory(categoryID uint) (count int64, err error) {
	err = DB.Model(&model.Product{}).
		Where("category_id = ? AND on_sale = ?", categoryID, true).
		Count(&count).Error
	return
}
