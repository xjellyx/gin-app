package usecase

import (
	"context"
	"gin-app/internal/domain"
	"gin-app/internal/domain/types"
	"github.com/casbin/casbin/v2"
	"github.com/spf13/cast"
	"strconv"
	"time"

	"gin-app/internal/domain/request"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"gorm.io/gorm/clause"
)

type casbinRuleUsecase struct {
	cfg CasbinRuleUsecaseConfig
}

type CasbinRuleUsecaseConfig struct {
	SysApiRepo domain.SysAPIRepo
	MenuRepo   domain.MenuRepo
	Casbin     casbin.IEnforcer
	Timeout    time.Duration
}

func NewCasbinRuleUsecase(cfg CasbinRuleUsecaseConfig) domain.CasbinRuleUsecase {
	return &casbinRuleUsecase{
		cfg: cfg,
	}

}

// AddAPIPermission 添加API权限
func (m *casbinRuleUsecase) AddAPIPermission(ctx context.Context, t types.CasbinRuleKey, req *request.AddAPIPermissionReq) error {
	var policies [][]string
	for _, v := range req.APIIds {
		var val int
		switch v.(type) {
		case int:
			val = v.(int)
		case float64:
			val = int(v.(float64))
		default:
			continue
		}
		api, err := m.cfg.SysApiRepo.FindOne(ctx, uint(val))
		if err != nil {
			return err
		}
		policies = append(policies, []string{string(t), strconv.Itoa(req.ID), api.Path, api.Method})
	}
	_, err := m.cfg.Casbin.RemoveFilteredPolicy(0, string(t), strconv.Itoa(req.ID))
	if err != nil {
		return serror.Error(serror.ErrCasbinRemoveFail, scontext.GetLanguage(ctx))
	}
	if len(policies) > 0 {
		_, err = m.cfg.Casbin.AddPolicies(policies)
		if err != nil {
			return serror.Error(serror.ErrCasbinAddFail, scontext.GetLanguage(ctx))
		}
	}
	return nil
}

// GetAPIPermission 获取API权限
func (m *casbinRuleUsecase) GetAPIPermission(ctx context.Context, t types.CasbinRuleKey, id uint) []uint {
	policies := m.cfg.Casbin.GetFilteredPolicy(0, string(t), strconv.Itoa(int(id)))
	var ids = make([]uint, 0)
	for _, v := range policies {
		var cond []clause.Expression
		cond = append(cond, clause.Eq{Column: "path", Value: v[2]})
		cond = append(cond, clause.Eq{Column: "method", Value: v[3]})
		by, err := m.cfg.SysApiRepo.FindOneBy(ctx, cond)
		if err != nil {
			return nil
		}
		ids = append(ids, by.ID)
	}
	return ids
}

// AddRoleMenuPermission 添加角色菜单权限
func (m *casbinRuleUsecase) AddRoleMenuPermission(ctx context.Context, t types.CasbinRuleKey, req *request.AddRoleMenuPermissionReq) error {
	var (
		policies [][]string
	)
	for _, v := range req.MenuIds {
		_, err := m.cfg.MenuRepo.FindOne(ctx, uint(v))
		if err != nil {
			return err
		}

		policies = append(policies, []string{string(t), strconv.Itoa(req.ID), strconv.Itoa(v), " "})
	}
	_, err := m.cfg.Casbin.RemoveFilteredPolicy(0, string(t), strconv.Itoa(req.ID))
	if err != nil {
		return serror.Error(serror.ErrCasbinRemoveFail, scontext.GetLanguage(ctx))
	}
	if len(policies) > 0 {
		_, err = m.cfg.Casbin.AddPolicies(policies)
		if err != nil {
			return serror.Error(serror.ErrCasbinAddFail, scontext.GetLanguage(ctx))
		}
	}
	return err
}

// GetRoleMenuPermission 获取角色菜单权限
func (m *casbinRuleUsecase) GetRoleMenuPermission(ctx context.Context, t types.CasbinRuleKey, id uint) []uint {
	policies := m.cfg.Casbin.GetFilteredPolicy(0, string(t), strconv.Itoa(int(id)))
	var ids = make([]uint, 0)
	for _, v := range policies {
		ids = append(ids, cast.ToUint(v[2]))
	}
	return ids
}

func (m *casbinRuleUsecase) SetRoleRouteFrontPage(ctx context.Context, req *request.SetRoleRouteFrontPageReq) error {
	_, err := m.cfg.Casbin.RemoveFilteredPolicy(0, string(types.CasbinRoleRouteFrontPageKey), strconv.Itoa(req.ID))
	if err != nil {
		return serror.Error(serror.ErrCasbinRemoveFail, scontext.GetLanguage(ctx))
	}

	_, err = m.cfg.Casbin.AddPolicy([]string{string(types.CasbinRoleRouteFrontPageKey), strconv.Itoa(req.ID), req.RoutePath, " "})
	if err != nil {
		return serror.Error(serror.ErrCasbinAddFail, scontext.GetLanguage(ctx))
	}
	return nil
}

func (m *casbinRuleUsecase) GetRoleRouteFrontPage(ctx context.Context, id uint) string {
	policies := m.cfg.Casbin.GetFilteredPolicy(0, string(types.CasbinRoleRouteFrontPageKey), strconv.Itoa(int(id)))
	if len(policies) == 0 {
		return ""
	}
	return policies[0][2]
}
