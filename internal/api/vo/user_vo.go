package vo

import "xiaomi-mall/internal/model"

// UserRegisterResp 注册响应
type UserRegisterResp struct {
	UserID   uint   `json:"user_id"`
	Phone string `json:"phone"`
	NickName string `json:"nick_name"`
}

// UserLoginResp 登录响应
type UserLoginResp struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"user_info"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Phone string `json:"phone"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Money    int64  `json:"money"`
	Role     int  `json:"role"`
	Status   string `json:"status"`
}

// NewUserInfo 从 Model 构造 VO
func NewUserInfo(user *model.User) UserInfo {
	return UserInfo{
		UserID:   user.ID,
		Phone: user.Phone,
		NickName: user.NickName,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Money:    user.Money,
		Role:     user.Role,
		Status:   user.Status,
	}
}
