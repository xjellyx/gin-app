package domain

import (
	"context"
	"time"
)

type UserInfo struct {
	Username  string     `json:"username"`  // username
	Email     string     `json:"email"`     // 邮箱
	Phone     string     `json:"phone"`     // 电话
	Gender    UserGender `json:"gender"`    // 性别
	CreatedAt time.Time  `json:"createdAt"` // 注册时间
}

type UserHimSelfUsecase interface {
	Info(ctx context.Context) (*UserInfo, error)
}
