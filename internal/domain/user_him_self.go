package domain

import (
	"context"

	"gin-app/internal/domain/response"
)

type UserHimSelfUsecase interface {
	Info(ctx context.Context) (*response.UserInfo, error)
	GetMenusTree(ctx context.Context, code string) (*response.GetMenusTreeResp, error)
	GetUserRoles(ctx context.Context) ([]*response.Role, error)
}
