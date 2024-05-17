package request

import (
	"gin-app/internal/domain/types"
	"gorm.io/datatypes"
)

type AddMenuReq struct {
	ParentId        int                                   `json:"parentId" `                                // 父级ID
	MenuType        types.MenuType                        `json:"menuType" binding:"required,min=1,max=3"`  // 菜单类型
	MenuName        string                                `json:"menuName" binding:"required"`              // 菜单名
	RouteName       string                                `json:"routeName"`                                // 路由名
	RoutePath       string                                `json:"routePath"`                                // 路由路径
	Component       string                                `json:"component"`                                // 组件路径
	Order           int                                   `json:"order"`                                    // 排序
	I18NKey         string                                `json:"i18nKey"`                                  // 国际化key
	Icon            string                                `json:"icon"`                                     // 图标
	IconType        types.MenuIconType                    `json:"iconType" binding:"required,min=1,max=2" ` // 图标类型
	Status          types.MenuStatus                      `json:"status" binding:"required,min=1,max=3"`    // 状态
	HideInMenu      bool                                  `json:"hideInMenu"`                               // 是否隐藏
	KeepAlive       bool                                  `json:"keepAlive"`                                // 是否缓存
	Constant        bool                                  `json:"constant"`                                 // 是否常量
	Href            string                                `json:"href"`                                     // 跳转链接
	MultiTab        bool                                  `json:"multiTab"`                                 // 是否多标签
	FixedIndexInTab int                                   `json:"fixedIndexInTab"`                          // 固定在标签上的index
	Query           datatypes.JSONSlice[types.MenuQuery]  `json:"query"`                                    // 查询参数
	Buttons         datatypes.JSONSlice[types.MenuButton] `json:"buttons"`                                  // 按钮
}

type GetMenusReq struct {
	Page
}

type EditMenuReq struct {
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
}

type AddMenuAPIPermissionReq struct {
	MenuId int   `json:"menuId" binding:"required"` // 菜单ID
	APIIds []any `json:"apiIds"`                    // apiIds
}
