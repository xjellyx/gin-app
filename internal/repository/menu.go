package repository

import (
	"gin-app/internal/domain"

	gormgenerics "github.com/xjellyx/gorm-generics"
)

type menuRepo struct {
	gormgenerics.BasicRepo[domain.Menu]
	db gormgenerics.Database
}

var _ domain.MenuRepo = (*menuRepo)(nil)

func NewMenuRepo(data gormgenerics.Database) domain.MenuRepo {
	d := gormgenerics.NewBasicRepository[domain.Menu](data)
	stu := &menuRepo{
		d,
		data,
	}
	return stu
}
