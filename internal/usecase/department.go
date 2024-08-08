package usecase

import (
	"context"
	"errors"

	"gin-app/internal/domain"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DepartmentUsecaseCfg struct {
	DepartmentRepo domain.DepartmentRepo
}

type departmentUsecase struct {
	cfg *DepartmentUsecaseCfg
}

func NewDepartmentUsecase(cfg *DepartmentUsecaseCfg) domain.DepartmentUsecase {
	return &departmentUsecase{
		cfg: cfg,
	}
}

// GetDepartments 获取部门列表
func (d *departmentUsecase) GetDepartments(ctx context.Context, req *request.GetDepartmentsReq) (*response.GetDepartmentsResp, error) {
	var cond []clause.Expression
	if req.Name != "" {
		cond = append(cond, clause.Like{
			Column: "name",
			Value:  "%" + req.Name + "%",
		})
	}
	if req.Principal != "" {
		cond = append(cond, clause.Like{
			Column: "principal",
			Value:  "%" + req.Principal + "%",
		})
	}
	cond = append(cond, req.SetOrmExpression())
	data, err := d.cfg.DepartmentRepo.Find(ctx, cond...)
	if err != nil {
		return nil, err
	}
	ret := &response.GetDepartmentsResp{
		Departments: make([]*response.Department, 0),
	}
	for _, i := range data {
		ret.Departments = append(ret.Departments, domain.Department2Resp(i))
	}
	return ret, nil
}

func (d *departmentUsecase) GetAllDepartments(ctx context.Context) ([]*response.Department, error) {
	//TODO implement me
	panic("implement me")
}

// GetDepartment 获取部门
func (d *departmentUsecase) GetDepartment(ctx context.Context, id uint) (*response.Department, error) {
	data, err := d.cfg.DepartmentRepo.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, serror.Error(serror.ErrDepartmentNotExist, scontext.GetLanguage(ctx))
		}
		return nil, err
	}
	return domain.Department2Resp(data), nil

}

// AddDepartment 添加部门
func (d *departmentUsecase) AddDepartment(ctx context.Context, req *request.AddDepartmentReq) error {
	if req.Pid != 0 {
		// 判断上级部门是否存在
		count, err := d.cfg.DepartmentRepo.Count(ctx, clause.Eq{
			Column: "id",
			Value:  req.Pid,
		})
		if err != nil {
			return err
		}
		if count == 0 {
			return serror.Error(serror.ErrDepartmentNotExist, scontext.GetLanguage(ctx))
		}
	}
	// 判断是否存在同名部门
	count, err := d.cfg.DepartmentRepo.Count(ctx, clause.Eq{
		Column: "name",
		Value:  req.Name,
	},
		clause.Eq{
			Column: "pid",
			Value:  req.Pid,
		},
	)
	if err != nil {
		return err
	}

	if count > 0 {
		return serror.Error(serror.ErrDepartmentExist, scontext.GetLanguage(ctx))
	}
	// 新增部门
	model := domain.AddDepartmentReq2Model(req)
	if err = d.cfg.DepartmentRepo.Create(ctx, model); err != nil {
		return err
	}
	return nil

}

func (d *departmentUsecase) DeleteBatch(ctx context.Context, ids []uint) error {
	//TODO implement me
	panic("implement me")
}

// EditDepartment 编辑部门
func (d *departmentUsecase) EditDepartment(ctx context.Context, id uint, req *request.EditDepartmentReq) error {
	if req.Pid != 0 {
		// 判断上级部门是否存在
		count, err := d.cfg.DepartmentRepo.Count(ctx, clause.Eq{
			Column: "id",
			Value:  req.Pid,
		})
		if err != nil {
			return err
		}
		if count == 0 {
			return serror.Error(serror.ErrDepartmentNotExist, scontext.GetLanguage(ctx))
		}
	}
	// 判断是否存在同名部门
	count, err := d.cfg.DepartmentRepo.Count(ctx, clause.Eq{
		Column: "name",
		Value:  req.Name,
	},
		clause.Eq{
			Column: "pid",
			Value:  req.Pid,
		},
		clause.Neq{
			Column: "id",
			Value:  id,
		},
	)
	if err != nil {
		return err
	}

	if count > 0 {
		return serror.Error(serror.ErrDepartmentExist, scontext.GetLanguage(ctx))
	}
	model := domain.EditDepartmentReq2Model(req)
	if err := d.cfg.DepartmentRepo.Update(ctx, id, model); err != nil {
		return err
	}
	return nil
}

// Delete 删除部门,包括子部门
func (d *departmentUsecase) Delete(ctx context.Context, id uint) error {
	return d.cfg.DepartmentRepo.ExecTX(ctx, func(ctx context.Context) error {
		var fc func(ctx context.Context, id uint) error
		fc = func(ctx context.Context, id uint) error {
			var children []*domain.Department
			children, err := d.cfg.DepartmentRepo.Find(ctx, clause.Eq{
				Column: "pid",
				Value:  id,
			})
			if err != nil {
				return err
			}
			// 递归删除所有子节点
			if len(children) > 0 {
				for _, v := range children {
					if err := fc(ctx, v.ID); err != nil {
						return err
					}
				}
			}
			// 删除当前节点
			if err = d.cfg.DepartmentRepo.DeleteOne(ctx, id); err != nil {
				return err
			}
			return nil

		}
		return fc(ctx, id)
	})
}

// GetDepartmentTree 获取部门树
func (d *departmentUsecase) GetDepartmentTree(ctx context.Context) (*response.GetDepartmentTreeResp, error) {
	var fc func(ctx context.Context, id uint, tree *response.DepartmentTree) error
	fc = func(ctx context.Context, id uint, tree *response.DepartmentTree) error {
		children, err := d.cfg.DepartmentRepo.Find(ctx, clause.Eq{
			Column: "pid",
			Value:  id,
		})
		if err != nil {
			return err
		}
		// 递归获取子节点
		if len(children) > 0 {
			for _, v := range children {
				child := domain.Department2Tree(v)
				if err = fc(ctx, v.ID, child); err != nil {
					return err
				}
				tree.Children = append(tree.Children, child)
			}
		}
		return nil
	}
	tree := make([]*response.DepartmentTree, 0)
	t := &response.DepartmentTree{
		Children: make([]*response.DepartmentTree, 0),
	}
	if err := fc(ctx, 0, t); err != nil {
		return nil, err
	}

	if len(t.Children) > 0 {
		tree = t.Children
	}
	return &response.GetDepartmentTreeResp{
		Tree: tree,
	}, nil
}
