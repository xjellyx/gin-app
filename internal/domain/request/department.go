package request

import "gin-app/internal/domain/types"

// AddDepartmentReq 添加部门
type AddDepartmentReq struct {
	Pid                  uint                   `json:"pid"`                     // 上级部门
	Name                 string                 `json:"name" binding:"required"` // 部门名
	Principal            string                 `json:"principal"`               // 部门负责人
	PrincipalPhoneNumber string                 `json:"principalPhoneNumber"`    // 负责人电话
	PrincipalEmail       string                 `json:"principalEmail"`          // 负责人邮箱
	Sequence             int                    `json:"sequence"`                // 排序
	Remark               string                 `json:"remark"`                  // 部门备注
	Status               types.DepartmentStatus `json:"status"`                  // 部门状态 1:启用 2:禁用
}

// EditDepartmentReq 编辑部门
type EditDepartmentReq struct {
	Pid                  uint                   `json:"pid"`                  // 上级部门
	Name                 string                 `json:"name"`                 // 部门名
	Principal            string                 `json:"principal"`            // 部门负责人
	PrincipalPhoneNumber string                 `json:"principalPhoneNumber"` // 负责人电话
	PrincipalEmail       string                 `json:"principalEmail"`       // 负责人邮箱
	Sequence             int                    `json:"sequence"`             // 排序
	Remark               string                 `json:"remark"`               // 部门备注
	Status               types.DepartmentStatus `json:"status"`               // 部门状态 1:启用 2:禁用
}

type GetDepartmentsReq struct {
	Name      string `form:"name"`      // 部门名
	Principal string `form:"principal"` // 部门负责人
	Page
}
