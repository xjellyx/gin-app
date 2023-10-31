package bootstrap

import (
	"gin-app/internal/infra/cache"
	"gin-app/pkg/serror"
	"gin-app/pkg/slog"

	gormgenerics "github.com/olongfen/gorm-generics"
	"github.com/ulule/limiter/v3"
	"go.uber.org/zap"
)

type Application struct {
	Conf     *Conf
	Log      *zap.Logger
	Database gormgenerics.Database
	Rdb      cache.Cache
	Limiter  *limiter.Limiter
}

var GlobalLog *zap.Logger

func App(confPath string) (*Application, error) {
	logger := slog.NewProduceLogger()
	GlobalLog = logger
	// 初始化多语言错误
	if err := serror.InitI18n(); err != nil {
		return nil, err
	}
	conf, err := NewConf(confPath)
	if err != nil {
		return nil, err
	}
	database, err := NewDatabase(conf, logger)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.NewRDB(cache.Config{
		Addr:      conf.RDBAddr,
		DB:        conf.RDBDB,
		Password:  conf.RDBPassword,
		KeyPrefix: conf.RDBKeyPrefix,
	})
	if err != nil {
		return nil, err
	}
	limit, err := NewLimitRate(conf)
	if err != nil {
		return nil, err
	}
	app := &Application{Database: database, Conf: conf, Log: logger, Rdb: rdb, Limiter: limit}
	return app, nil
}

func (a *Application) Close() error {
	if a == nil {
		return nil
	}
	a.Database.Close()
	return nil
}
