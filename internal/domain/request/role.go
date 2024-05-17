package request

import "gin-app/internal/domain/types"

type GetRolesReq struct {
	Page
	Name   string           `form:"name"`   // 名称
	Code   string           `form:"code"`   // 编码
	Status types.RoleStatus `form:"status"` // 状态
}

type AddRoleReq struct {
	Name   string `json:"name"`   // 名称
	Code   string `json:"code"`   // 编码
	Desc   string `json:"desc"`   // 描述
	Status string `json:"status"` // 状态
}

type AddRoleAPIPermissionReq struct {
	RoleId int   `json:"roleId" binding:"required"` // 角色ID
	APIIds []any `json:"apiIds"`                    // apiIds
}
