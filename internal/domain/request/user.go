package request

import (
	"gin-app/internal/domain/types"
)

// UserAdminListReq 用户列表请求
type UserAdminListReq struct {
	Page
	Username string           `form:"username"` // 用户名
	Email    string           `form:"email"`    // 邮箱
	Phone    string           `form:"phone"`    // 手机号码
	Gender   types.UserGender `form:"gender"`   // 性别
	Status   types.UserStatus `form:"status"`   // 状态
}

// UserAdminAddReq 用户添加
type UserAdminAddReq struct {
	Username string           `json:"username"`                              // 用户名
	Email    string           `json:"email" binding:"required,email"`        // 邮箱
	Phone    string           `json:"phone"`                                 // 手机号码
	Gender   types.UserGender `json:"gender" binding:"required,min=1,max=3"` // 性别1未知，2男，3女
	Status   types.UserStatus `json:"status" binding:"required,min=1,max=3"` // 状态：1正常，2禁止，3删除
	Roles    []string         `json:"roles"`                                 // 角色
}

// SignupReq 用户注册
type SignupReq struct {
	//Username     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required_without=Email"`
	Email    string `binding:"required_without=Phone" json:"email"`
	Password string `json:"password" binding:"required,min=8,max=16"`
	Code     string `json:"code"` // 验证码
}

// SignInReq 用户登录
type SignInReq struct {
	Phone    string `json:"phone" binding:"required_without=Email"`
	Email    string `binding:"required_without=Phone" json:"email"`
	Password string `json:"password" binding:"required,min=8,max=16"`
}

type UserAdminUpdateReq struct {
	Username string           `json:"username"` // 用户名
	Email    string           `json:"email"`    // 邮箱
	Phone    string           `json:"phone"`    // 手机号码
	Gender   types.UserGender `json:"gender"`   // 性别
	Status   types.UserStatus `json:"status"`   // 状态
	Roles    []string         `json:"roles"`    // 角色
}
