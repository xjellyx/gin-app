package orm

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// _tableNamePrefix 表格头部，需要自己定义
	_tableNamePrefix = ""
)

// DriverName 数据库驱动名称
type DriverName string

const (
	Postgresql DriverName = "postgresql"
	MySQL                 = "mysql"
	SQLite                = "sqlite"
	SQLServer             = "sqlserver"
	TiDB       DriverName = "tidb"
)

type Modeler interface {
	TableName() string
}

// DbConnect init database connect
func DbConnect(driver DriverName, dsn string, opts ...Option) (db *gorm.DB, err error) {
	var (
		gormConfig = &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: _tableNamePrefix,
			},
		}
		option = options{
			logger: zap.NewNop(),
		}
	)
	for _, o := range opts {
		o.apply(&option)
	}
	//
	switch driver {
	case Postgresql:
		//dbCfg := c.DB
		//dns := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s`, dbCfg.IP,
		//	dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Port, dbCfg.SSLMode, dbCfg.TimeZone)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			return
		}
	case MySQL:
		db, err = gorm.Open(mysql.New(mysql.Config{
			DriverName: "go-sql-driver",
			DSN:        dsn,
		}), gormConfig)
		if err != nil {
			return
		}
	case SQLite:
		db, err = gorm.Open(sqlite.Open(dsn), gormConfig)
		if err != nil {
			return
		}
	}

	if db == nil {
		err = errors.New("database connect failed")
		return
	}
	// true 自动迁移
	if option.autoMigrate && len(option.autoMigrateDst) > 0 {
		if err = db.AutoMigrate(option.autoMigrateDst...); err != nil {
			return
		}
	}
	// 使用链路追踪
	if option.opentracingPlugin != nil {
		if err = db.Use(option.opentracingPlugin); err != nil {
			return
		}
	}
	return db, nil
}
