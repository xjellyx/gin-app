package bootstrap

import (
	"gin-app/internal/domain"
	"gin-app/pkg/serror"

	"go.uber.org/zap"
)

type Application struct {
	Conf     *Conf
	Log      *zap.Logger
	Database domain.Database
}

func App(confPath string) (*Application, error) {
	// 初始化多语言错误
	if err := serror.InitI18n(); err != nil {
		return nil, err
	}
	conf, err := NewConf(confPath)
	if err != nil {
		return nil, err
	}
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
