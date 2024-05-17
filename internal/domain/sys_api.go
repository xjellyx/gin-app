package domain

import (
	"context"
	"gin-app/internal/domain/response"
	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/gorm"
)

type SysAPI struct {
	gorm.Model
	Path    string `gorm:"size:255;uniqueIndex:idx_sys_apis_path_method;comment:api路径"`
	Method  string `gorm:"size:10;uniqueIndex:idx_sys_apis_path_method;comment:http请求方法"`
	Summary string `gorm:"size:1024;comment:说明"`
	Tag     string `gorm:"size:1024;comment:tag"`
	Title   string `gorm:"size:256;comment:名称信息"`
}

type SysAPIRepo interface {
	gormgenerics.BasicRepo[SysAPI]
}

type SysAPIUsecase interface {
	GetSysAPITree(ctx context.Context) (response.SysAPITreeRespList, error)
}
