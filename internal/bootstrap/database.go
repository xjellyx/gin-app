package bootstrap

import (
	"gin-app/internal/domain"
	"gin-app/pkg/orm"

	"go.uber.org/zap"
)

// NewDatabase 新建数据库
func NewDatabase(conf *Conf, logger *zap.Logger) (domain.Database, error) {
	dataBase, err := orm.NewDataBase(conf.DBDriver, conf.DBDsn,
		orm.WithAutoMigrate(conf.DBAutoMigrate),
		orm.WithAutoMigrateDst([]any{&domain.User{}}),
		orm.WithLogger(logger),
		orm.WithOpentracingPlugin(&orm.OpentracingPlugin{}),
	)
	if err != nil {
		return nil, err
	}
	return dataBase, nil
}
