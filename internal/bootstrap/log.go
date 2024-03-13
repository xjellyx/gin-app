package bootstrap

import "gopkg.in/natefinch/lumberjack.v2"

func NewLumberjack(configs ...LumberjackConfig) *lumberjack.Logger {
	var (
		conf LumberjackConfig
	)
	if len(configs) != 0 {
		conf = configs[0]
	}
	if conf.LogPath == "" {
		conf.LogPath = "./log/log.log"
	}
	if conf.MaxSize == 0 {
		conf.MaxSize = 200
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = 60
	}
	if conf.MaxBackups == 0 {
		conf.MaxBackups = 50
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.LogPath,
		MaxSize:    conf.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: conf.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     conf.MaxAge,     // 保留旧文件的最大天数
		Compress:   conf.Compress,   // 是否压缩/归档旧文件
	}
	return lumberJackLogger

}
