package domain

import (
	"context"
	"gin-app/internal/domain/response"
)

type UserHimSelfUsecase interface {
	Info(ctx context.Context) (*response.UserInfo, error)
}
