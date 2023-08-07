package cache

import (
	"context"
	"time"
)

// Cache 令牌存储接口
type Cache interface {
	// Set 存储令牌数据，并指定到期时间
	Set(ctx context.Context, uuid string, val string, expiration time.Duration) error
	Get(ctx context.Context, uuid string) (string, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, uuid string) (bool, error)
	Del(ctx context.Context, key string) error
	DelByKeyPrefix(ctx context.Context, keyPrefix string) error
	Exists(ctx context.Context, key string) bool
	// Close 关闭存储
	Close() error
}
