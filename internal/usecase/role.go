package usecase

import (
	"context"
	"gin-app/internal/domain"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm/clause"
	"time"
)

type roleUsecase struct {
	cfg RoleUsecaseConfig
}

type RoleUsecaseConfig struct {
	Repo       domain.RoleRepo
	SysAPIRepo domain.SysAPIRepo
	Casbin     casbin.IEnforcer
	MenuRepo   domain.MenuRepo
	Timeout    time.Duration
}

func NewRoleUsecase(cfg RoleUsecaseConfig) domain.RoleUsecase {
	return &roleUsecase{
		cfg: cfg,
	}
}

// GetRoles 获取角色列表
func (r *roleUsecase) GetRoles(ctx context.Context, req *request.GetRolesReq) (*response.GetRolesResp, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()
	var cond []clause.Expression
	if req.Name != "" {
		cond = append(cond, clause.Like{Column: "name", Value: "%" + req.Name + "%"})
	}
	if req.Code != "" {
		cond = append(cond, clause.Like{Column: "code", Value: "%" + req.Code + "%"})
	}
	if req.Status != 0 {
		cond = append(cond, clause.Eq{Column: "status", Value: req.Status})
	}
	count, err := r.cfg.Repo.Count(ctx, cond...)
	if err != nil {
		return nil, err
	}
	cond = append(cond, req.SetOrmExpression())
	// 从数据库获取数据
	roles, err := r.cfg.Repo.Find(ctx, cond...)
	if err != nil {
		return nil, err
	}

	ret := &response.GetRolesResp{
		Records: make([]*response.Role, 0),
		Pagination: response.Pagination{
			Total:    count,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
	}
	for _, i := range roles {
		d := domain.Role2Resp(i)
		ret.Records = append(ret.Records, d)
	}

	return ret, nil
}

// GetAllRoles 获取全部角色
func (r *roleUsecase) GetAllRoles(ctx context.Context) ([]*response.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()
	roles, err := r.cfg.Repo.Find(ctx, []clause.Expression{clause.Eq{Column: "status", Value: types.RoleStatusEnable}}...)
	if err != nil {
		return nil, err
	}
	ret := make([]*response.Role, 0)
	for _, i := range roles {
		d := domain.Role2Resp(i)
		ret = append(ret, d)
	}
	return ret, nil
}

// AddRole 添加角色
func (r *roleUsecase) AddRole(ctx context.Context, req *request.AddRoleReq) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()
	role := &domain.Role{
		Name:   req.Name,
		Code:   req.Code,
		Desc:   req.Desc,
		Status: types.RoleStatusEnable,
	}
	err := r.cfg.Repo.Create(ctx, role)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除用户
func (r *roleUsecase) Delete(ctx context.Context, id uint) error {
	ctx, cancelFunc := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancelFunc()

	return r.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		err := r.cfg.Repo.DeleteOne(ctx, id)
		if err != nil {
			return err
		}

		// 删除关联的用户角色
		err = r.cfg.Repo.Database().DB(ctx).Table("user_roles").Where("role_id = ?", id).Delete(ctx).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *roleUsecase) DeleteBatch(ctx context.Context, ids []uint) error {
	ctx, cancelFunc := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancelFunc()
	cond := clause.IN{}
	for _, v := range ids {
		cond.Values = append(cond.Values, v)
	}
	cond.Column = "id"
	err := r.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		err := r.cfg.Repo.Database().DB(ctx).Table("user_roles").Clauses(clause.IN{Column: "role_id", Values: cond.Values}).Delete(ctx).Error
		if err != nil {
			return err
		}
		return r.cfg.Repo.DeleteBy(ctx, []clause.Expression{cond})
	})
	if err != nil {
		return err
	}
	return nil
}
