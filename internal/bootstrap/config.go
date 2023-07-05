package bootstrap

import (
	"log"
	"strings"

	"gin-app/pkg/orm"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Conf 配置环境
type Conf struct {
	HTTPort       uint           `mapstructure:"HTTP_PORT"`
	DBDriver      orm.DriverName `mapstructure:"DB_DRIVER"`
	DBDsn         string         `mapstructure:"DB_DSN"`
	DBAutoMigrate bool           `mapstructure:"DB_AUTO_MIGRATE"`
}

func NewConf(configPath string) (*Conf, error) {
	// 设置默认值
	viper.SetDefault("HTTP_PORT", 8080)

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	//
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	// 读取终端变量
	_ = viper.BindPFlags(pflag.CommandLine)
	conf := &Conf{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file:%s Op:%s\n", e.Name, e.Op)
		if err := viper.Unmarshal(conf); err != nil {
			log.Fatal(err)
		}
	})
	return conf, nil
}
