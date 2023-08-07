package cache

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Config struct {
	Addr      string
	DB        int
	Password  string
	KeyPrefix string
}

// NewRDB 创建基于redis存储实例
func NewRDB(cfg Config) (Cache, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RDB{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}, nil
}

// RDB redis存储
type RDB struct {
	cli    *redis.Client
	prefix string
}

func (s *RDB) wrapperKey(key string) string {
	return fmt.Sprintf("%s|%s", s.prefix, key)
}

func (s *RDB) Get(ctx context.Context, key string) (string, error) {
	return s.cli.Get(ctx, s.wrapperKey(key)).Result()
}

// Set ...
func (s *RDB) Set(ctx context.Context, uuid string, val string, expiration time.Duration) error {
	cmd := s.cli.Set(ctx, s.wrapperKey(uuid), val, expiration)
	return cmd.Err()
}

// Delete ...
func (s *RDB) Delete(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Del(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// Check ...
func (s *RDB) Check(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Exists(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *RDB) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.cli.SetNX(ctx, s.wrapperKey(key), value, expiration).Err()
}

func (s *RDB) Del(ctx context.Context, key string) error {
	return s.cli.Del(ctx, s.wrapperKey(key)).Err()
}

func (s *RDB) DelByKeyPrefix(ctx context.Context, keyPrefix string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = s.cli.Scan(ctx, cursor, keyPrefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err = s.cli.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *RDB) Exists(ctx context.Context, key string) bool {
	return s.cli.Exists(ctx, s.wrapperKey(key)).Val() == 1
}

// Close ...
func (s *RDB) Close() error {
	return s.cli.Close()
}
