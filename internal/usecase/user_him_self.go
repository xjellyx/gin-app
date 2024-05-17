package usecase

import (
	"context"
	"fmt"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"
	"github.com/casbin/casbin/v2"
	"strconv"
	"time"

	"gin-app/internal/domain"
	"gin-app/pkg/scontext"

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
		fmt.Println("2222", v)
		ret.Routes = append(ret.Routes, domain.MenuModel2UserMenuResp(v))
	}
	ret.Routes = response.BuildMenuTree(ret.Routes)
	fmt.Println("sssssss", ret.Routes[0])
	return ret, nil
}
