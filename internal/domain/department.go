package domain

import (
	"context"

	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"

	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/gorm"
)

// Department 部门表模型
type Department struct {
	gorm.Model
	Pid                  uint                   `gorm:"index;default:0;unique:idx_department_pid_name;comment:上级部门"`
	Name                 string                 `gorm:"size:256;index;unique:idx_department_pid_name;comment:部门名"`
	Principal            string                 `gorm:"default:null;comment:部门负责人"`
	PrincipalPhoneNumber string                 `gorm:"default:null;comment:负责人电话"`
	PrincipalEmail       string                 `gorm:"default:null;comment:负责人邮箱"`
	Sequence             int                    `gorm:"default:null;comment:排序"`
	Remark               string                 `gorm:"default:null;comment:部门备注"`
	Status               types.DepartmentStatus `gorm:"index;comment:部门状态,1:启用,2:禁用"`
}

// DepartmentRepo 部门仓库
type DepartmentRepo interface {
	gormgenerics.BasicRepo[Department]
}

// DepartmentUsecase 部门业务
type DepartmentUsecase interface {
	GetDepartments(ctx context.Context, req *request.GetDepartmentsReq) (*response.GetDepartmentsResp, error)
	GetAllDepartments(ctx context.Context) ([]*response.Department, error)
	GetDepartment(ctx context.Context, id uint) (*response.Department, error)
	GetDepartmentTree(ctx context.Context) (*response.GetDepartmentTreeResp, error)
	AddDepartment(ctx context.Context, req *request.AddDepartmentReq) error
	DeleteBatch(ctx context.Context, ids []uint) error
	EditDepartment(ctx context.Context, id uint, req *request.EditDepartmentReq) error
	Delete(ctx context.Context, id uint) error
}

func Department2Resp(i *Department) *response.Department {
	return &response.Department{
		ID:                   i.ID,
		Pid:                  i.Pid,
		Name:                 i.Name,
		Remark:               i.Remark,
		Status:               i.Status,
		Principal:            i.Principal,
		PrincipalPhoneNumber: i.PrincipalPhoneNumber,
		PrincipalEmail:       i.PrincipalEmail,
		Sequence:             i.Sequence,
	}
}

func AddDepartmentReq2Model(req *request.AddDepartmentReq) *Department {
	return &Department{
		Pid:                  req.Pid,
		Name:                 req.Name,
		Principal:            req.Principal,
		PrincipalPhoneNumber: req.PrincipalPhoneNumber,
		PrincipalEmail:       req.PrincipalEmail,
		Sequence:             req.Sequence,
		Remark:               req.Remark,
		Status:               req.Status,
	}
}

func EditDepartmentReq2Model(req *request.EditDepartmentReq) *Department {
	return &Department{
		Pid:                  req.Pid,
		Name:                 req.Name,
		Principal:            req.Principal,
		PrincipalPhoneNumber: req.PrincipalPhoneNumber,
		PrincipalEmail:       req.PrincipalEmail,
		Sequence:             req.Sequence,
		Remark:               req.Remark,
		Status:               req.Status,
	}

}

func Department2Tree(i *Department) *response.DepartmentTree {
	return &response.DepartmentTree{
		Department: *Department2Resp(i),
		Children:   make([]*response.DepartmentTree, 0),
	}
}
