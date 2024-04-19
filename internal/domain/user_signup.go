package domain

import (
	"context"
	"time"
)

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

// SignInResp 用户登录响应
type SignInResp struct {
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	RefreshToken string    `json:"refreshToken"`
}

// SignupUsecase 用户注册用例
type SignupUsecase interface {
	Signup(ctx context.Context, req *SignupReq) error
	SignIn(ctx context.Context, req *SignInReq) (*SignInResp, error)
}
