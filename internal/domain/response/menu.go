package response

import (
	"gin-app/internal/domain/types"
)

type GetMenusResp struct {
	Menus      []*Menu    `json:"records"`
	Pagination Pagination `json:"pagination"`
}

type Menu struct {
	ID              int                `json:"id"`
	ParentId        int                `json:"parentId" `                                // 父级ID
	MenuType        types.MenuType     `json:"menuType" binding:"required,min=1,max=3"`  // 菜单类型
	MenuName        string             `json:"menuName" binding:"required"`              // 菜单名
	RouteName       string             `json:"routeName"`                                // 路由名
	RoutePath       string             `json:"routePath"`                                // 路由路径
	Component       string             `json:"component"`                                // 组件路径
	Order           int                `json:"order"`                                    // 排序
	I18NKey         string             `json:"i18nKey"`                                  // 国际化key
	Icon            string             `json:"icon"`                                     // 图标
	IconType        types.MenuIconType `json:"iconType" binding:"required,min=1,max=2" ` // 图标类型
	Status          types.MenuStatus   `json:"status" binding:"required,min=1,max=3"`    // 状态
	HideInMenu      bool               `json:"hideInMenu"`                               // 是否隐藏
	KeepAlive       bool               `json:"keepAlive"`                                // 是否缓存
	Constant        bool               `json:"constant"`                                 // 是否常量
	Href            string             `json:"href"`                                     // 跳转链接
	MultiTab        bool               `json:"multiTab"`                                 // 是否多标签
	FixedIndexInTab int                `json:"fixedIndexInTab"`                          // 固定在标签上的index
	Query           []types.MenuQuery  `json:"query"`                                    // 查询参数
	Buttons         []types.MenuButton `json:"buttons"`                                  // 按钮
	Children        []*Menu            `json:"children"`
}

type MenuTreeItem struct {
	ID       int             `json:"id"`
	Label    string          `json:"label"`
	PId      int             `json:"pId"`
	Order    int             `json:"order"`
	I18NKey  string          `json:"i18nKey"` // 国际化key
	Children []*MenuTreeItem `json:"children"`
}
