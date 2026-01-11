package dto

// UserRegisterReq 注册请求
type UserRegisterReq struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	NickName string `json:"nick_name" binding:"required,min=2,max=30"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
}

// UserLoginReq 登录请求
type UserLoginReq struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Password string `json:"password" binding:"required"`
}
