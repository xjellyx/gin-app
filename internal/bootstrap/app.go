package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"gin-app/internal/infra/cache"
	"gin-app/pkg/serror"

	"github.com/casbin/casbin/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lmittmann/tint"
	"github.com/ulule/limiter/v3"
	gormgenerics "github.com/xjellyx/gorm-generics"
)

type Application struct {
	Conf     *Conf
	Database gormgenerics.Database
	Rdb      cache.Cache
	Limiter  *limiter.Limiter
	Casbin   casbin.IEnforcer
}

var GlobalApp *Application

func App(confPath string) (*Application, error) {

	// 初始化多语言错误
	if err := serror.InitI18n(); err != nil {
		return nil, err
	}
	conf, err := NewConf(confPath)
	if err != nil {
		return nil, err
	}
	if !conf.LogConf.IsProd {
		w := os.Stderr
		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.DateTime,
			}),
		))
	} else {
		logger := slog.NewJSONHandler(NewLumberjack(conf.LogConf), nil)
		slog.SetDefault(slog.New(logger))
	}

	database, err := NewDatabase(conf)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.NewRDB(cache.Config{
		Addr:      conf.RDB.Addr,
		DB:        conf.RDB.DB,
		Password:  conf.RDB.Password,
		KeyPrefix: conf.RDB.Prefix,
	})
	if err != nil {
		return nil, err
	}
	limit, err := NewLimitRate(conf)
	if err != nil {
		return nil, err
	}
	en, err := NewEnforcer(database, conf.CasbinModel)
	if err != nil {
		return nil, err
	}
	app := &Application{Database: database, Conf: conf, Rdb: rdb, Limiter: limit}
	app.Casbin = en
	GlobalApp = app
	return app, nil
}

func (a *Application) Close() error {
	if a == nil {
		return nil
	}
	a.Database.Close()
	return nil
}
