package bootstrap

import (
	redis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewLimitRate(cfg *Conf) (*limiter.Limiter, error) {
	// 每分钟限制一百次请求
	rate, err := limiter.NewRateFromFormatted("100-M")
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.RDB.Addr,
		DB:       cfg.RDB.DB + 1,
		Password: cfg.RDB.Password,
	})
	store, err := sredis.NewStoreWithOptions(cli, limiter.StoreOptions{
		Prefix:   "limiter_gin_example",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, err
	}
	return limiter.New(store, rate), nil
}
