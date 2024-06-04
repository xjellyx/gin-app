package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gin-app/docs"
	"gin-app/internal/domain"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"
	"gin-app/pkg/sslog"

	"github.com/jackc/pgx/v5/pgconn"
	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/gorm"
)

// NewDatabase 新建数据库
func NewDatabase(conf *Conf) (gormgenerics.Database, error) {
	var logger *slog.Logger
	if conf.LogConf.IsProd {
		c := conf.LogConf
		c.LogPath = "log/db.log"
		l := slog.NewJSONHandler(NewLumberjack(c), nil)
		logger = slog.New(l)
	} else {
		logger = slog.Default()
	}
	dataBase, err := gormgenerics.NewDataBase(gormgenerics.DriverName(conf.DB.Driver), conf.DB.Dsn,
		gormgenerics.WithAutoMigrate(conf.DB.AutoMigrate),
		gormgenerics.WithAutoMigrateDst([]any{&domain.User{}, &domain.Role{}, &domain.Menu{}, &domain.SysAPI{}}),
		gormgenerics.WithLogger(sslog.NewDBLog(logger)),
		gormgenerics.WithOpentracingPlugin(&gormgenerics.OpentracingPlugin{}),
		gormgenerics.WithTranslateError(translateErr),
		gormgenerics.WithTablePrefix(conf.DB.Prefix),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		err = autoCreateRestfulAPIByDoc(dataBase.DB(context.Background()))
		if err != nil {
			slog.Error("autoCreateRestfulAPIByDoc", "msg", err)
		}
	}()
	return dataBase, nil
}

func translateErr(ctx context.Context, db *gorm.DB) (err error) {
	err = db.Error
	lan := scontext.GetLanguage(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		switch db.Statement.Table {
		case "users":
			return serror.Error(serror.ErrUserRecordNotFound, lan)
		default:
			return serror.Error(serror.ErrRecordNotFound, lan)
		}
	}

	var errV *pgconn.PgError
	ok := errors.As(err, &errV)
	if !ok {
		return
	}
	switch errV.Code {
	case "23505":
		switch errV.TableName {
		case "users":
			return domain.TranslateUserDBErr(errV, lan)
		}
	}

	return err
}

func autoCreateRestfulAPIByDoc(db *gorm.DB) error {
	var (
		methods = []string{
			"GET",
			"HEAD",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
		}
	)
	m := docs.GetAPIData()
	for k, v := range m["paths"].(map[string]any) {
		for _, method := range methods {
			if val, ok := v.(map[string]any)[strings.ToLower(method)]; ok {
				if val.(map[string]any)["summary"] == nil {
					fmt.Println(v)
				}
				for _, tag := range val.(map[string]any)["tags"].([]any) {
					data := &domain.SysAPI{Summary: val.(map[string]any)["summary"].(string)}
					data.Tag = tag.(string)
					if data.Tag == "UserHimSelf" || data.Tag == "UserSign 用户注册登录" {
						data.White = true
					}

					cond := &domain.SysAPI{Path: k, Method: method, Tag: data.Tag}
					var count int64
					db.Model(&domain.SysAPI{}).Where(cond).Count(&count)
					if count == 1 {
						data.Path = k
						data.Method = method
						if err := db.Model(&domain.SysAPI{}).Where(cond).Updates(data).Error; err != nil {
							return err
						}
					} else {
						if err := db.Where(cond).Attrs(data).FirstOrCreate(
							&domain.SysAPI{}).Error; err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
