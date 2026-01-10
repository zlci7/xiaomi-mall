// service/user_service.go
package service

import (
	"time"
	"xiaomi-mall/config"
	"xiaomi-mall/internal/api/dto"
	"xiaomi-mall/internal/api/vo"
	"xiaomi-mall/internal/dao"
	"xiaomi-mall/internal/model"
	"xiaomi-mall/pkg/encrypt"
	"xiaomi-mall/pkg/jwtx"
	"xiaomi-mall/pkg/xerr"
)

type UserService struct{}

var User = new(UserService)

// Register 接收 DTO，返回 VO
func (s *UserService) Register(req dto.UserRegisterReq) (*vo.UserRegisterResp, error) {
	// 1. 业务校验
	exist, err := dao.User.ExistOrNotByUserName(req.UserName)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}
	if exist {
		return nil, xerr.NewErrCode(xerr.USER_ALREADY_EXISTS)
	}

	// 2. 密码加密（业务逻辑）
	passwordDigest, err := encrypt.EncryptPassword(req.Password)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.USER_ENCRYPT_ERROR)
	}

	// 3️⃣ DTO → Model 转换（在 Service 层完成）
	userModel := &model.User{
		UserName:       req.UserName,
		Email:          req.Email,
		PasswordDigest: passwordDigest,
		NickName:       req.NickName,
		Avatar:         req.Avatar,
		Status:         "active",
		Money:          0,
		Role:           0,
	}

	// 4. 调用 DAO 保存（传 Model）
	if err := dao.User.CreateUser(userModel); err != nil {
		return nil, xerr.NewErrCode(xerr.USER_CREATE_ERROR)
	}

	// 5️⃣ Model → VO 转换（在 Service 层完成）
	resp := &vo.UserRegisterResp{
		UserID:   userModel.ID,
		UserName: userModel.UserName,
		NickName: userModel.NickName,
	}

	// 6️⃣ 返回 VO（不是 Model）
	return resp, nil
}

// Login 接收 DTO，返回 VO
func (s *UserService) Login(req dto.UserLoginReq) (*vo.UserLoginResp, error) {
	// 1. 查找用户（DAO 返回 Model）
	user, err := dao.User.GetUserByUserName(req.UserName)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.USER_NOT_FOUND)
	}

	// 2. 密码校验
	if !encrypt.ValidatePassword(req.Password, user.PasswordDigest) {
		return nil, xerr.NewErrCode(xerr.USER_PASSWORD_ERROR)
	}

	// 3. 生成 Token
	token, err := jwtx.GetToken(config.AppConfig.Jwt.AccessSecret, time.Now().Unix(), config.AppConfig.Jwt.AccessExpire, int64(user.ID))

	if err != nil {
		return nil, xerr.NewErrCode(xerr.TOKEN_GEN_ERROR)
	}

	// 4️⃣ Model → VO 转换
	resp := &vo.UserLoginResp{
		Token:    token,
		UserInfo: vo.NewUserInfo(user), // 使用 VO 的构造函数
	}

	// 5️⃣ 返回 VO
	return resp, nil
}
