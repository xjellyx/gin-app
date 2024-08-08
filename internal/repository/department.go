package repository

import (
	"gin-app/internal/domain"

	gormgenerics "github.com/xjellyx/gorm-generics"
)

type departmentRepo struct {
	gormgenerics.BasicRepo[domain.Department]
	db gormgenerics.Database
}

func NewDepartmentRepo(db gormgenerics.Database) domain.DepartmentRepo {
	return &departmentRepo{
		BasicRepo: gormgenerics.NewBasicRepository[domain.Department](db),
		db:        db,
	}
}
