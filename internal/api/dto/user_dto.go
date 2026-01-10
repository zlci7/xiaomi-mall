package dto

// UserRegisterReq 注册请求
type UserRegisterReq struct {
	UserName string `json:"user_name" binding:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=40"`
	NickName string `json:"nick_name" binding:"required,min=2,max=30"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
}

// UserLoginReq 登录请求
type UserLoginReq struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
