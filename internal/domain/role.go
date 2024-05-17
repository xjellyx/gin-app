package domain

import (
	"context"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"
	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name   string           `gorm:"size:256;index;comment:角色名"`
	Code   string           `gorm:"size:256;uniqueIndex;comment:角色编码"`
	Desc   string           `gorm:"comment:角色描述"`
	Status types.RoleStatus `gorm:"index;comment:角色状态,1:启用,2:禁用,3:删除"`
}

type RoleRepo interface {
	gormgenerics.BasicRepo[Role]
}

type RoleUsecase interface {
	GetRoles(ctx context.Context, req *request.GetRolesReq) (*response.GetRolesResp, error)
	GetAllRoles(ctx context.Context) ([]*response.Role, error)
	AddRole(ctx context.Context, req *request.AddRoleReq) error
	DeleteBatch(ctx context.Context, ids []uint) error
	Delete(ctx context.Context, id uint) error
}

func Role2Resp(i *Role) *response.Role {
	return &response.Role{
		ID:     i.ID,
		Name:   i.Name,
		Code:   i.Code,
		Desc:   i.Desc,
		Status: i.Status,
	}

}
