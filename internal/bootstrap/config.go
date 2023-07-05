package bootstrap

import "gin-app/pkg/orm"

// Conf 配置环境
type Conf struct {
	DBDriver      orm.DriverName `mapstructure:"DB_DRIVER"`
	DBDsn         string         `mapstructure:"DB_DSN"`
	DBAutoMigrate bool           `mapstructure:"DB_AUTO_MIGRATE"`
}

func NewConf(configPath string) *Conf {
	// TODO
	conf := &Conf{}
	return conf
}
