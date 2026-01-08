package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	gorm.Model
	UserName       string `gorm:"unique" json:"user_name"`
	Email          string `json:"email"`
	PasswordDigest string `json:"-"` // 密码不返回给前端
	NickName       string `json:"nick_name"`
	Avatar         string `json:"avatar"`
	Status         string `json:"status"`                 // Active, Suspended
	Money          int64  `gorm:"default:0" json:"money"` // 余额，单位：分
	Role           int32  `gorm:"default:0" json:"role"`  // 角色：0:普通用户 1:管理员 2:超级管理员
}

// Address 收货地址表
type Address struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"` // 加索引方便查询
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	IsDefault bool   `gorm:"default:false" json:"is_default"` // 是否默认地址
}
