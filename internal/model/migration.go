package model

import "gorm.io/gorm"

// Migrate 自动迁移数据库结构
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Address{},
		&Category{},
		&Product{},
		&ProductSku{},
		&Order{},
		&OrderItem{},
		&SeckillProduct{},
		&SeckillOrder{},
	)
	return err
}
