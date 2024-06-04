package usecase

import (
	"context"
	"strconv"
	"time"

	"gin-app/internal/domain"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"
	"gin-app/pkg/scontext"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm/clause"
)

type UserHimSelfUsecaseConfig struct {
	Repo           domain.UserRepo
	Casbin         casbin.IEnforcer
	RoleRepo       domain.RoleRepo
	MenuRepo       domain.MenuRepo
	ContextTimeout time.Duration
}

type userHimSelfUsecase struct {
	cfg UserHimSelfUsecaseConfig
}

func NewUserHimSelfUsecase(cfg UserHimSelfUsecaseConfig) domain.UserHimSelfUsecase {
	return &userHimSelfUsecase{cfg: cfg}
}

func (u *userHimSelfUsecase) Info(ctx context.Context) (*response.UserInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancel()
	uid := scontext.GetUserUuid(ctx)
	cond := clause.Eq{
		Column: "uuid",
		Value:  uid,
	}
	ret, err := u.cfg.Repo.FindOneBy(ctx, []clause.Expression{cond})
	if err != nil {
		return nil, err
	}
	data := &response.UserInfo{
		Username:  ret.Username,
		Email:     ret.Email,
		Phone:     ret.Phone,
		Gender:    ret.Gender,
		CreatedAt: ret.CreatedAt,
		Roles:     scontext.GetRoles(ctx),
	}
	return data, nil
}

// GetMenusTree get menus tree
func (u *userHimSelfUsecase) GetMenusTree(ctx context.Context, code string) (*response.GetMenusTreeResp, error) {
	roles := scontext.GetRoles(ctx)
	ret := &response.GetMenusTreeResp{
		Routes: make([]*response.UserMenuItem, 0),
		Home:   "",
	}
	if len(roles) == 0 {
		return ret, nil
	}
	hasRole := false
	for _, v := range roles {
		if v == code {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return ret, nil
	}
	role, err := u.cfg.RoleRepo.FindOneBy(ctx, []clause.Expression{clause.Eq{Column: "code", Value: code}})
	if err != nil {
		return nil, err
	}
	policiesHomePage := u.cfg.Casbin.GetFilteredPolicy(0, string(types.CasbinRoleRouteFrontPageKey), strconv.Itoa(int(role.ID)))
	if len(policiesHomePage) > 0 {
		ret.Home = policiesHomePage[0][2]
	}

	policies := u.cfg.Casbin.GetFilteredPolicy(0, string(types.CasbinRoleMenuKey), strconv.Itoa(int(role.ID)))
	if len(policies) == 0 {
		return ret, nil
	}
	condIn := clause.IN{
		Column: "id",
		Values: nil,
	}
	for _, v := range policies {
		condIn.Values = append(condIn.Values, v[2])
	}
	menus, err := u.cfg.MenuRepo.Find(ctx, condIn, clause.Eq{Column: "constant", Value: false})
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return ret, nil
	}
	for _, v := range menus {
		ret.Routes = append(ret.Routes, domain.MenuModel2UserMenuResp(v))
	}
	ret.Routes = response.BuildMenuTree(ret.Routes)
	return ret, nil
}

func (u *userHimSelfUsecase) GetUserRoles(ctx context.Context) ([]*response.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancel()
	rolesCode := scontext.GetRoles(ctx)
	if len(rolesCode) == 0 {
		return make([]*response.Role, 0), nil
	}
	cond := clause.IN{
		Column: "code",
		Values: nil,
	}
	for _, v := range rolesCode {
		cond.Values = append(cond.Values, v)
	}
	roles, err := u.cfg.RoleRepo.Find(ctx, []clause.Expression{cond,
		clause.Eq{Column: "status", Value: types.RoleStatusEnable}}...)
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
