package repository

import (
	"context"
	"errors"
	"gin-app/internal/domain"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// basicRepo 基础存储库实现
type basicRepo[T any] struct {
	database domain.Database
}

// NewBasicRepo 新建基础存储库
func NewBasicRepo[T any](database domain.Database) domain.BasicRepo[T] {
	return &basicRepo[T]{database}
}

// Create 创建一条记录
func (b *basicRepo[T]) Create(ctx context.Context, ent *T) error {
	db := b.database.DB(ctx).Create(ent)
	if err := db.Error; err != nil {
		return translateErr(ctx, db)
	}
	return nil
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
		return
	}

	errV, ok := err.(*pgconn.PgError)
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

func processCond(db *gorm.DB, conds []clause.Expression) *gorm.DB {
	for _, v := range conds {
		val, ok := v.(clause.OrderBy)
		if ok {
			for _, order := range val.Columns {
				db = db.Order(order)
			}

			continue
		}
		db = db.Where(v)
	}
	return db
}

// Find 查询
func (b *basicRepo[T]) Find(ctx context.Context, limit domain.Limit, conds ...clause.Expression) ([]*T, int64, error) {
	var (
		model T
	)
	db := b.database.DB(ctx).Model(&model).Debug()
	db = processCond(db, conds)
	var (
		count int64
	)
	// 如果需要全部数据
	if limit.Count {
		if err := db.Count(&count).Error; err != nil {
			err = translateErr(ctx, db)
			return nil, 0, err
		}
	}
	var (
		data []*T
	)
	if !limit.All {
		switch {
		case limit.PageSize > 0 && limit.PageNum > 0:
			db = db.Offset(int((limit.PageNum - 1) * limit.PageSize)).Limit(int(limit.PageSize))
		case limit.PageSize > 0 && limit.PageNum == 0:
			db = db.Limit(int(limit.PageSize))
		default:
			limit.PageSize = 10
			limit.PageNum = 1
			db = db.Offset(int((limit.PageNum - 1) * limit.PageSize)).Limit(int(limit.PageSize))
		}
	}

	if err := db.Find(&data).Error; err != nil {
		err = translateErr(ctx, db)
		return nil, 0, err
	}
	return data, count, nil
}

// Delete 删除一条数据
func (b *basicRepo[T]) Delete(ctx context.Context, id int) error {
	var (
		model T
	)
	db := b.database.DB(ctx).Where("id = ?", id).Delete(&model)
	if err := db.Error; err != nil {
		return translateErr(ctx, db)
	}
	return nil
}

// FindOne 查询一条数据
func (b *basicRepo[T]) FindOne(ctx context.Context, id int) (*T, error) {
	var (
		model T
		data  *T
	)
	db := b.database.DB(ctx).Model(&model).Where("id = ?", id).First(&data)
	if err := db.Error; err != nil {
		err = translateErr(ctx, db)
		return nil, err
	}
	return data, nil
}

// FindOneBy 多条件查询一条数据
func (b *basicRepo[T]) FindOneBy(ctx context.Context, conds ...clause.Expression) (*T, error) {
	var (
		model T
		data  *T
	)
	db := b.database.DB(ctx).Model(&model)
	db = processCond(db, conds)
	if err := db.First(&data).Error; err != nil {
		err = translateErr(ctx, db)
		return nil, err
	}
	return data, nil
}

// Update 更新一条数据
func (b *basicRepo[T]) Update(ctx context.Context, id int, ent *T) error {
	var (
		model T
	)
	db := b.database.DB(ctx).Model(&model).Where("id = ?", id).Updates(ent)
	if err := db.Error; err != nil {
		err = translateErr(ctx, db)
		return err
	}
	return nil
}
