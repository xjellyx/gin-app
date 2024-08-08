package response

import "gin-app/internal/domain/types"

type Department struct {
	ID                   uint                   `json:"id"`                   // 部门ID
	Pid                  uint                   `json:"pid"`                  // 上级部门
	Name                 string                 `json:"name"`                 // 部门名
	Remark               string                 `json:"remark"`               // 部门备注
	Status               types.DepartmentStatus `json:"status"`               // 部门状态
	Principal            string                 `json:"principal"`            // 部门负责人
	PrincipalPhoneNumber string                 `json:"principalPhoneNumber"` // 负责人电话
	PrincipalEmail       string                 `json:"principalEmail"`       // 负责人邮箱
	Sequence             int                    `json:"sequence"`             // 排序
}

// GetDepartmentsResp 部门列表
type GetDepartmentsResp struct {
	Departments []*Department `json:"departments"`
}

// GetDepartmentTreeResp 部门树
type GetDepartmentTreeResp struct {
	Tree []*DepartmentTree `json:"tree"`
}

// DepartmentTree 部门树
type DepartmentTree struct {
	Department
	Children []*DepartmentTree `json:"children"`
}
