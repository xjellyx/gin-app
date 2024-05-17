package repository

import (
	"gin-app/internal/domain"
	gormgenerics "github.com/xjellyx/gorm-generics"
)

type sysAPIRepo struct {
	gormgenerics.BasicRepo[domain.SysAPI]
	data gormgenerics.Database
}

func NewSysAPIRepo(data gormgenerics.Database) domain.SysAPIRepo {
	return &sysAPIRepo{data: data, BasicRepo: gormgenerics.NewBasicRepository[domain.SysAPI](data)}
}
