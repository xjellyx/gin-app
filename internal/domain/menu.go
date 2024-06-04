package domain

import (
	"context"

	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"

	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	CreateBy        string                                `gorm:"size:256;comment:创建人"`
	UpdateBy        string                                `gorm:"size:256;comment:更新人"`
	Status          types.MenuStatus                      `gorm:"index;comment:状态,1:启用,2:禁用,3:删除"`
	ParentId        int                                   `gorm:"index;comment:父级ID"`
	MenuType        types.MenuType                        `gorm:"index;comment:菜单类型,1:目录,2:菜单,3:按钮"`
	MenuName        string                                `gorm:"size:256;comment:菜单名"`
	RouteName       string                                `gorm:"size:256;comment:路由名"`
	RoutePath       string                                `gorm:"size:256;comment:路由路径"`
	Component       string                                `gorm:"size:256;comment:组件路径"`
	Order           int                                   `gorm:"index;comment:排序"`
	I18NKey         string                                `gorm:"size:256;comment:国际化key"`
	Icon            string                                `gorm:"size:256;comment:图标"`
	IconType        types.MenuIconType                    `gorm:"index;comment:图标类型"`
	HideInMenu      bool                                  `gorm:"default:false;comment:是否隐藏"`
	KeepAlive       bool                                  `gorm:"default:false;comment:是否缓存"`
	Constant        bool                                  `gorm:"default:false;comment:是否常量"`
	Href            string                                `gorm:"comment:跳转链接"`
	MultiTab        bool                                  `gorm:"default:false;comment:是否多标签"`
	FixedIndexInTab int                                   `gorm:"comment:固定在标签上的index"`
	Query           datatypes.JSONSlice[types.MenuQuery]  `gorm:"comment:查询参数"`
	Buttons         datatypes.JSONSlice[types.MenuButton] `gorm:"comment:按钮"`
}

type MenuRepo interface {
	gormgenerics.BasicRepo[Menu]
}

type MenuUsecase interface {
	AddMenu(ctx context.Context, req *request.AddMenuReq) error
	GetMenus(ctx context.Context, req *request.GetMenusReq) (*response.GetMenusResp, error)
	DeleteMenu(ctx context.Context, menuId uint) error
	DeleteMenus(ctx context.Context, menuId []uint) error
	EditMenu(ctx context.Context, id uint, req *request.EditMenuReq) error
	GetMenuTree(ctx context.Context) ([]*response.MenuTreeItem, error)
	GetAllPages(ctx context.Context) ([]string, error)
	CheckRouteExist(ctx context.Context, routeName string) bool
}

func MenuForm2Model(req *request.AddMenuReq) *Menu {
	return &Menu{
		ParentId:        req.ParentId,
		MenuType:        req.MenuType,
		MenuName:        req.MenuName,
		RouteName:       req.RouteName,
		RoutePath:       req.RoutePath,
		Component:       req.Component,
		Order:           req.Order,
		I18NKey:         req.I18NKey,
		Icon:            req.Icon,
		IconType:        req.IconType,
		Status:          req.Status,
		HideInMenu:      req.HideInMenu,
		KeepAlive:       req.KeepAlive,
		Constant:        req.Constant,
		Href:            req.Href,
		MultiTab:        req.MultiTab,
		FixedIndexInTab: req.FixedIndexInTab,
		Query:           req.Query,
		Buttons:         req.Buttons,
	}
}

func MenuModel2Form(menu *Menu) *response.Menu {
	return &response.Menu{
		ID:              int(menu.ID),
		ParentId:        menu.ParentId,
		MenuType:        menu.MenuType,
		MenuName:        menu.MenuName,
		RouteName:       menu.RouteName,
		RoutePath:       menu.RoutePath,
		Component:       menu.Component,
		Order:           menu.Order,
		I18NKey:         menu.I18NKey,
		Icon:            menu.Icon,
		IconType:        menu.IconType,
		Status:          menu.Status,
		HideInMenu:      menu.HideInMenu,
		KeepAlive:       menu.KeepAlive,
		Constant:        menu.Constant,
		Href:            menu.Href,
		MultiTab:        menu.MultiTab,
		FixedIndexInTab: menu.FixedIndexInTab,
		Query:           menu.Query,
		Buttons:         menu.Buttons,
	}
}

func MenuModel2TreeForm(m *Menu) *response.MenuTreeItem {
	return &response.MenuTreeItem{
		ID:       int(m.ID),
		Label:    m.MenuName,
		PId:      m.ParentId,
		Order:    m.Order,
		I18NKey:  m.I18NKey,
		Children: make([]*response.MenuTreeItem, 0),
	}

}

func MenuEdit2Model(req *request.EditMenuReq) *Menu {
	return &Menu{
		MenuType:        req.MenuType,
		MenuName:        req.MenuName,
		RouteName:       req.RouteName,
		RoutePath:       req.RoutePath,
		Component:       req.Component,
		Order:           req.Order,
		I18NKey:         req.I18NKey,
		Icon:            req.Icon,
		IconType:        req.IconType,
		Status:          req.Status,
		HideInMenu:      req.HideInMenu,
		KeepAlive:       req.KeepAlive,
		Constant:        req.Constant,
		Href:            req.Href,
		MultiTab:        req.MultiTab,
		FixedIndexInTab: req.FixedIndexInTab,
		Query:           req.Query,
		Buttons:         req.Buttons,
	}
}

func MenuModel2UserMenuResp(m *Menu) *response.UserMenuItem {
	return &response.UserMenuItem{
		ID:        int(m.ID),
		PId:       m.ParentId,
		Name:      m.RouteName,
		Path:      m.RoutePath,
		Component: m.Component,
		Meta: response.UserMenuItemMeta{
			Title:       m.MenuName,
			Icon:        m.Icon,
			KeepAlive:   m.KeepAlive,
			DefaultMenu: m.HideInMenu,
			CloseTab:    m.KeepAlive,
			I18nKey:     m.I18NKey,
			Order:       m.Order,
		},
		Children: make([]*response.UserMenuItem, 0),
	}
}
