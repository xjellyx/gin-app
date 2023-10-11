package bootstrap

import (
	"context"
	"errors"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"
	"gin-app/pkg/slog"

	"gin-app/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	gormgenerics "github.com/olongfen/gorm-generics"
	"github.com/olongfen/gorm-generics/achieve"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewDatabase 新建数据库
func NewDatabase(conf *Conf, logger *zap.Logger) (gormgenerics.Database, error) {
	dataBase, err := achieve.NewDataBase(conf.DBDriver, conf.DBDsn,
		achieve.WithAutoMigrate(conf.DBAutoMigrate),
		achieve.WithAutoMigrateDst([]any{&domain.User{}}),
		achieve.WithLogger(slog.NewDBLog(logger)),
		achieve.WithOpentracingPlugin(&achieve.OpentracingPlugin{}),
	)
	if err != nil {
		return nil, err
	}
	return dataBase, nil
}

func translateErr(ctx context.Context, db *gorm.DB) (err error) {
	err = db.Error
	lan := scontext.GetLanguage(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		switch db.Statement.Table {
		case domain.User{}.TableName():
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
		case domain.User{}.TableName():
			return domain.TranslateUserDBErr(errV, lan)
		}
	}

	return err
}
