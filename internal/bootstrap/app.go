package bootstrap

import (
	"gin-app/internal/domain"
	"gin-app/pkg/orm"
	"gin-app/pkg/selferror"

	"go.uber.org/zap"
)

type Application struct {
	Conf     *Conf
	Log      *zap.Logger
	Database domain.Database
}

func App(confPath string) (*Application, error) {
	// 初始化多语言错误
	if err := selferror.InitI18n(); err != nil {
		return nil, err
	}
	conf := NewConf(confPath)
	conf.DBAutoMigrate = true
	conf.DBDriver = orm.Postgresql
	conf.DBDsn = "postgres://postgres:123456@192.168.3.4:5432/public?sslmode=disable"
	database, err := NewDatabase(conf)
	if err != nil {
		return nil, err
	}
	app := &Application{Database: database, Conf: conf}
	return app, nil
}

func (a *Application) Close() error {
	if a == nil {
		return nil
	}
	a.Database.Close()
	return nil
}
