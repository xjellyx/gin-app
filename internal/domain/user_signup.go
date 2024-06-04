package domain

import (
	"context"

	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
)

// SignupUsecase 用户注册用例
type SignupUsecase interface {
	Signup(ctx context.Context, req *request.SignupReq) error
	SignIn(ctx context.Context, req *request.SignInReq) (*response.SignInResp, error)
	GetConstantMenuTree(ctx context.Context) ([]*response.UserMenuItem, error)
}
