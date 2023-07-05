package domain

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BasicRepo 基本的存储库
type BasicRepo[T any] interface {
	// Create 创建一条数据
	Create(ctx context.Context, ent *T) error
	// Find 查询
	Find(ctx context.Context, limit Limit, conds ...clause.Expression) ([]*T, int64, error)
	Update(ctx context.Context, id int, ent *T) error
	FindOneBy(ctx context.Context, conds ...clause.Expression) (*T, error)
	FindOne(ctx context.Context, id int) (*T, error)
	Delete(ctx context.Context, id int) error
}

type Database interface {
	DB(ctx context.Context) *gorm.DB
	Close() error
	ExecTx(ctx context.Context, fc func(context.Context) error) error
}

// Limit 数据库数量查询限制
type Limit struct {
	All      bool // true获取全部条目
	PageSize uint // 每页数量
	PageNum  uint // 页数
	Count    bool // true获取总数
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
