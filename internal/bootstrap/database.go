package bootstrap

import (
	"gin-app/internal/domain"
	"gin-app/pkg/orm"
	"gin-app/pkg/slog"
)

// NewDatabase 新建数据库
func NewDatabase(conf *Conf) (domain.Database, error) {
	dataBase, err := orm.NewDataBase(conf.DBDriver, conf.DBDsn,
		orm.WithAutoMigrate(conf.DBAutoMigrate),
		orm.WithAutoMigrateDst([]any{&domain.User{}}),
		orm.WithLogger(slog.NewProduceLogger()),
	)
	if err != nil {
		return nil, err
	}
	return dataBase, nil
}
