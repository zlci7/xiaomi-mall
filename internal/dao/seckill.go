package dao

import "xiaomi-mall/internal/model"

type SeckillDao struct{}

var Seckill = new(SeckillDao)

// ==================== 管理端：秒杀商品管理 ====================

// 创建秒杀商品入库
func (d *SeckillDao) CreateSeckillProduct(product *model.SeckillProduct) error {
	return DB.Create(product).Error

}

// 删除秒杀商品
func (d *SeckillDao) DeleteSeckillProduct(id uint) error {
	// return DB.Model(&model.SeckillProduct{}).Where("id =?", id).Delete(&model.SeckillProduct{}).Error
	return DB.Delete(&model.SeckillProduct{}, "id = ?", id).Error
}

// 手动开启/结束秒杀
func (d *SeckillDao) UpdateSeckillStatus(id uint, status int) error {
	return DB.Model(&model.SeckillProduct{}).Where("id=?", id).Update("status", status).Error
}

// ==================== 用户端：秒杀商品查询 ====================

// ==================== 用户端：秒杀商品下单 ====================
