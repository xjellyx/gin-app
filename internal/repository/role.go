package repository

import (
	"gin-app/internal/domain"

	gormgenerics "github.com/xjellyx/gorm-generics"
)

type roleRepo struct {
	gormgenerics.BasicRepo[domain.Role]
	data gormgenerics.Database
}

func NewRoleRepo(data gormgenerics.Database) domain.RoleRepo {
	return &roleRepo{data: data, BasicRepo: gormgenerics.NewBasicRepository[domain.Role](data)}
}
