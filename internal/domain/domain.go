package domain

import (
	"context"
	"gorm.io/gorm"
)

// BasicRepo 基本的存储库
type BasicRepo[T any] interface {
	// Create 创建一条数据
	Create(ctx context.Context, ent *T) error
	// Find
}

type Database interface {
	DB(ctx context.Context) *gorm.DB
	Close() error
	ExecTx(ctx context.Context, fc func(context.Context) error) error
}

// Response 返回体
type Response struct {
	// Code 0为成功，其他都是失败
	Code string `json:"code"`
	// Msg 返回提示
	Msg string `json:"msg"`
	// Data 返回数据
	Data any `json:"data"`
}
