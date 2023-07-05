package repository

import (
	"context"
	"gin-app/internal/domain"
)

// basicRepo 基础存储库实现
type basicRepo[T any] struct {
	database domain.Database
}

// NewBasicRepo 新建基础存储库
func NewBasicRepo[T any](database domain.Database) domain.BasicRepo[T] {
	return &basicRepo[T]{database}
}

// Create 创建一条记录
func (b *basicRepo[T]) Create(ctx context.Context, ent *T) error {
	if err := b.database.DB(ctx).Create(ent).Error; err != nil {
		return err
	}
	return nil
}
