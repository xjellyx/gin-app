package sslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBLog struct {
	slog                      *slog.Logger
	LogLevel                  slog.Level
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func (l *DBLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = slog.Level(level)
	return &newlogger
}

func (l *DBLog) Info(ctx context.Context, s string, i ...interface{}) {
	l.slog.InfoContext(ctx, s, i)
}

func (l *DBLog) Warn(ctx context.Context, s string, i ...interface{}) {
	l.slog.WarnContext(ctx, s, i)
}

func (l *DBLog) Error(ctx context.Context, s string, i ...interface{}) {
	l.slog.ErrorContext(ctx, s, i)
}

func (l *DBLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= slog.LevelDebug {
		return
	}

	elapsed := time.Since(begin)
	consuming := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	//line := utils.FileWithLineNum()
	//lineMsg := slog.String("Line", line)
	consumingMsg := slog.String("TimeConsuming", consuming)
	sql, rows := fc()
	sqlMsg := slog.String("SQL", sql)
	rowsMsg := slog.Int64("RowsAffected", rows)
	switch {
	case err != nil && l.LogLevel >= slog.LevelError && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		if rows == -1 {
			l.slog.ErrorContext(ctx, "DBTrace", "Error", err, consumingMsg, sqlMsg)
		} else {
			l.slog.ErrorContext(ctx, "DBTrace", "Error", err, consumingMsg, sqlMsg, rowsMsg)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= slog.LevelInfo:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.slog.WarnContext(ctx, "DBTrace", consumingMsg, "SlowLog", slowLog, sqlMsg)
		} else {
			l.slog.WarnContext(ctx, "DBTrace", consumingMsg, "SlowLog", slowLog, sqlMsg, rowsMsg)
		}
	case l.LogLevel == slog.LevelInfo:
		if rows == -1 {
			l.slog.InfoContext(ctx, "DBTrace", consumingMsg, sqlMsg)
		} else {
			l.slog.InfoContext(ctx, "DBTrace", consumingMsg, sqlMsg, rowsMsg)
		}
	}
}

func NewDBLog(logger *slog.Logger) logger.Interface {
	return &DBLog{
		slog:                      logger,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  slog.LevelInfo,
	}
}
