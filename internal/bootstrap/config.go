package bootstrap

import (
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/olongfen/gorm-generics/achieve"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Conf 配置环境
type Conf struct {
	HTTPort uint `mapstructure:"HTTP_PORT"`
	//
	DBDriver      achieve.DriverName `mapstructure:"DB_DRIVER"`
	DBDsn         string             `mapstructure:"DB_DSN"`
	DBAutoMigrate bool               `mapstructure:"DB_AUTO_MIGRATE"`
	//
	RDBAddr      string `mapstructure:"RDB_ADDR"`
	RDBPassword  string `mapstructure:"RDB_PASSWORD"`
	RDBDB        int    `mapstructure:"RDB_DB"`
	RDBKeyPrefix string `mapstructure:"RDB_KEY_PREFIX"`
	//
	JWTExpireTime        time.Duration `mapstructure:"JWT_EXPIRE_TIME"`
	JWTRefreshExpireTime time.Duration `mapstructure:"JWT_REFRESH_EXPIRE_TIME"`
	JWTSigningMethod     string        `mapstructure:"JWT_SIGNING_METHOD"`
	JWTSigningKey        string        `mapstructure:"JWT_SIGNING_KEY"`
	JWTRefreshSingingKey string        `mapstructure:"JWT_REFRESH_SIGNING_KEY"`
	JWTEnabled           bool          `mapstructure:"JWT_ENABLED"`
	// 日志
	LogConf LumberjackConfig
}

type LumberjackConfig struct {
	IsProd     bool
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
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
	GlobalConf = conf
	return conf, nil
}

func GetConfig() Conf {
	return *GlobalConf
}

var GlobalConf *Conf
