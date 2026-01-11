package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	gorm.Model
	// 核心认证信息
	Phone          string `gorm:"unique;not null" json:"phone"` // 手机号，唯一且必填
	PasswordDigest string `json:"-"`                            // 密码密文
	
	// 用户基本信息
	NickName string `json:"nick_name"` // 昵称 (显示用)
	Email    string `json:"email"`     // 邮箱 (可选，用于找回密码)
	Avatar   string `json:"avatar"`    // 头像 URL
	
	// 状态与权限
	Status string `gorm:"default:'Active'" json:"status"` // Active:正常, Suspended:封禁
	Role   int    `gorm:"default:0" json:"role"`          // 0:普通用户 1:管理员
	Money  int64  `gorm:"default:0" json:"money"`         // 余额 (单位:分，防止精度丢失)
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
