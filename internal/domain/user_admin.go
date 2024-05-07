package domain

import (
	"context"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
)

type UserAdminUsecase interface {
	List(ctx context.Context, req *request.UserAdminListReq) (*response.UserAdminListResp, error)
	Add(ctx context.Context, req *request.UserAdminAddReq) error
}
