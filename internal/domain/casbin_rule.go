package domain

import (
	"context"

	"gin-app/internal/domain/request"
	"gin-app/internal/domain/types"
)

type CasbinRule struct {
	ID    uint   `gorm:"primarykey"`
	Ptype string `gorm:"size:128;default:p;uniqueIndex:unique_index"`
	V0    string `gorm:"default:null;size:128;column:v0;uniqueIndex:unique_index"`
	V1    string `gorm:"default:null;size:128;column:v1;uniqueIndex:unique_index"`
	V2    string `gorm:"default:null;size:128;column:v2;uniqueIndex:unique_index"`
	V3    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V4    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V5    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
}

// CasbinRuleUsecase casbin rule usecase
type CasbinRuleUsecase interface {
	AddAPIPermission(ctx context.Context, t types.CasbinRuleKey, req *request.AddAPIPermissionReq) error
	GetAPIPermission(ctx context.Context, t types.CasbinRuleKey, id uint) []uint
	AddRoleMenuPermission(ctx context.Context, t types.CasbinRuleKey, req *request.AddRoleMenuPermissionReq) error
	GetRoleMenuPermission(ctx context.Context, t types.CasbinRuleKey, id uint) []uint
	SetRoleRouteFrontPage(ctx context.Context, req *request.SetRoleRouteFrontPageReq) error
	GetRoleRouteFrontPage(ctx context.Context, id uint) string
}
