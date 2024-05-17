package response

import "gin-app/internal/domain/types"

type GetRolesResp struct {
	Records    []*Role    `json:"records"`
	Pagination Pagination `json:"pagination"`
}

type Role struct {
	ID     uint             `json:"id"`
	Name   string           `json:"name"`   // 名称
	Code   string           `json:"code"`   // 编码
	Desc   string           `json:"desc"`   // 描述
	Status types.RoleStatus `json:"status"` // 状态
}
