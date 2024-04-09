package bootstrap

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Conf 配置环境
type Conf struct {
	HTTPort uint `mapstructure:"HTTP_PORT"`
	//
	DB  DBCfg  `mapstructure:"db"`
	RDB RDBCfg `mapstructrue:"rdb"`
	//
	JWT JWTCfg `mapstructure:"jwt"`
	//
	Resource string `mapstructure:"RESOURCE"`
	// 日志
	LogConf LumberjackConfig
	Nacos   NacosCfg `mapstructure:"NACOS"`
}

type DBCfg struct {
	Driver      string `mapstructure:"driver"`
	Dsn         string `mapstructure:"dsn"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
	Prefix      string `mapstructure:"prefix"`
}

type RDBCfg struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
}

type JWTCfg struct {
	ExpireTime        uint   `mapstructure:"expire_time"`
	RefreshExpireTime uint   `mapstructure:"refresh_expire_time"`
	SigningMethod     string `mapstructure:"signing_method"`
	SigningKey        string `mapstructure:"signing_key"`
	RefreshSingingKey string `mapstructure:"refresh_singing_key"`
	Enabled           bool   `mapstructure:"enabled"`
}

type LumberjackConfig struct {
	IsProd     bool
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type NacosCfg struct {
	IP           string  `mapstructure:"ip"`
	Port         uint64  `mapstructure:"port"`
	ClientName   string  `mapstructure:"client_name"`
	ClientIP     string  `mapstructure:"client_ip"`
	ClientPort   uint64  `mapstructure:"client_port"`
	ClientWeight float64 `mapstructure:"client_weight"`
}

func NewConf(configPath string) (*Conf, error) {
	// 设置默认值
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("RESOURCE", "resource")

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
	GlobalConf = conf
	fmt.Println("aaaaaaaa", conf.Nacos, conf.DB)
	return conf, nil
}

func GetConfig() Conf {
	return *GlobalConf
}

var GlobalConf *Conf
