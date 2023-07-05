package domain

import "context"

// SignupReq 用户注册
type SignupReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `binding:"required,email" json:"email"`
	Password string `json:"password" binding:"required"`
}

// SignupUsecase 用户注册用例
type SignupUsecase interface {
	Signup(ctx context.Context, req *SignupReq) error
}
