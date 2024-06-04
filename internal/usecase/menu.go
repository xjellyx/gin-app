package usecase

import (
	"context"
	"time"

	"gin-app/internal/domain"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm/clause"
)

type menuUsecase struct {
	cfg MenuUsecaseConfig
}

type MenuUsecaseConfig struct {
	Repo       domain.MenuRepo
	SysApiRepo domain.SysAPIRepo
	Casbin     casbin.IEnforcer
	Timeout    time.Duration
}

func NewMenuUsecase(cfg MenuUsecaseConfig) domain.MenuUsecase {
	return &menuUsecase{
		cfg: cfg,
	}
}

func (m *menuUsecase) AddMenu(ctx context.Context, req *request.AddMenuReq) error {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.Timeout)
	defer cancel()
	return m.cfg.Repo.Create(ctx, domain.MenuForm2Model(req))
}

func (m *menuUsecase) GetMenus(ctx context.Context, req *request.GetMenusReq) (*response.GetMenusResp, error) {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.Timeout)
	defer cancel()
	var cond []clause.Expression
	cond = append(cond, clause.Eq{Column: "parent_id", Value: 0})
	count, err := m.cfg.Repo.Count(ctx, cond...)
	if err != nil {
		return nil, err
	}

	cond = append(cond, clause.OrderBy{
		Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "order"}, Desc: false}},
	})
	cond = append(cond, req.SetOrmExpression())
	menus, err := m.cfg.Repo.Find(ctx, cond...)
	if err != nil {
		return nil, err
	}
	ret := &response.GetMenusResp{
		Menus:      make([]*response.Menu, 0),
		Pagination: response.Pagination{PageNum: req.PageNum, PageSize: req.PageSize, Total: count},
	}
	for _, v := range menus {
		d := domain.MenuModel2Form(v)
		d.Children, _ = m.getChildren(ctx, v.ID)
		ret.Menus = append(ret.Menus, d)
	}
	return ret, nil
}

func (m *menuUsecase) getChildren(ctx context.Context, parentId uint) ([]*response.Menu, error) {
	menus, err := m.cfg.Repo.Find(ctx, clause.Eq{Column: "parent_id", Value: parentId},
		clause.OrderBy{
			Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "order"}, Desc: false}},
		})
	if err != nil {
		return nil, err
	}
	ret := make([]*response.Menu, 0)
	for _, v := range menus {
		d := domain.MenuModel2Form(v)
		d.Children, _ = m.getChildren(ctx, v.ID)
		ret = append(ret, d)
	}
	return ret, nil
}

func (m *menuUsecase) DeleteMenu(ctx context.Context, menuId uint) error {
	return m.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		err := m.cfg.Repo.DeleteOne(ctx, menuId)
		if err != nil {
			return err
		}
		err = m.cfg.Repo.DeleteBy(ctx, []clause.Expression{clause.Eq{Column: "parent_id", Value: menuId}})
		if err != nil {
			return err
		}
		// todo 删除关联的权限信息
		return nil
	})
}

func (m *menuUsecase) DeleteMenus(ctx context.Context, menuIds []uint) error {
	return m.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		var val []any
		for _, v := range menuIds {
			val = append(val, v)
		}
		err := m.cfg.Repo.DeleteBy(ctx, []clause.Expression{clause.Eq{Column: "id", Value: val}})
		if err != nil {
			return err
		}
		err = m.cfg.Repo.DeleteBy(ctx, []clause.Expression{clause.Eq{Column: "parent_id", Value: val}})
		if err != nil {
			return err
		}
		// todo 删除关联的权限信息
		return nil
	})
}

func (m *menuUsecase) EditMenu(ctx context.Context, id uint, req *request.EditMenuReq) error {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.Timeout)
	defer cancel()
	data := domain.MenuEdit2Model(req)
	err := m.cfg.Repo.Update(ctx, id, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *menuUsecase) GetMenuTree(ctx context.Context) ([]*response.MenuTreeItem, error) {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.Timeout)
	defer cancel()
	menus, err := m.getTreeChildren(ctx, 0)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *menuUsecase) getTreeChildren(ctx context.Context, parentId uint) ([]*response.MenuTreeItem, error) {
	menus, err := m.cfg.Repo.Find(ctx, clause.Eq{Column: "parent_id", Value: parentId},
		clause.Eq{Column: "constant", Value: false},
		clause.OrderBy{
			Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "order"}, Desc: false}},
		})
	if err != nil {
		return nil, err
	}
	ret := make([]*response.MenuTreeItem, 0)
	for _, v := range menus {
		d := domain.MenuModel2TreeForm(v)
		d.Children, _ = m.getTreeChildren(ctx, v.ID)
		ret = append(ret, d)
	}
	return ret, nil
}

func (m *menuUsecase) GetAllPages(ctx context.Context) ([]string, error) {
	var (
		pages []string
	)
	err := m.cfg.Repo.Database().DB(ctx).Model(&domain.Menu{}).Select("route_name").Scan(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (m *menuUsecase) CheckRouteExist(ctx context.Context, routeName string) bool {
	_, err := m.cfg.Repo.FindOneBy(ctx, []clause.Expression{clause.Eq{Column: "route_name", Value: routeName}})
	if err != nil {
		return false
	}
	return true
}
