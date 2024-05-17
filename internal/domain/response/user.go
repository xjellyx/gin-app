package response

import (
	"sort"
	"time"

	"gin-app/internal/domain/types"
)

type UserInfo struct {
	Username  string           `json:"username"`  // username
	Email     string           `json:"email"`     // 邮箱
	Phone     string           `json:"phone"`     // 电话
	Gender    types.UserGender `json:"gender"`    // 性别 1未知、2男性、3女性
	CreatedAt time.Time        `json:"createdAt"` // 注册时间
	Roles     []string         `json:"roles"`     // 角色
}

type UserListInfo struct {
	Uuid   string           `json:"uuid"`
	ID     uint             `json:"id"`
	Status types.UserStatus `json:"status"`
	UserInfo
}

// UserAdminListResp 用户列表响应
type UserAdminListResp struct {
	Records    []*UserListInfo `json:"records"`
	Pagination *Pagination     `json:"pagination"`
}

// SignInResp 用户登录响应
type SignInResp struct {
	Token        string    `json:"token"`
	ExpireAt     time.Time `json:"expireAt"`
	RefreshToken string    `json:"refreshToken"`
}

type GetMenusTreeResp struct {
	Routes []*UserMenuItem `json:"routes"`
	Home   string          `json:"home"`
}

type UserMenuItem struct {
	ID        int              `json:"id"`
	PId       int              `json:"pId"`
	Name      string           `json:"name"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Meta      UserMenuItemMeta `json:"meta"`
	Children  []*UserMenuItem  `json:"children"`
}
type UserMenuItemMeta struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	KeepAlive   bool   `json:"keepAlive"`
	DefaultMenu bool   `json:"defaultMenu"`
	CloseTab    bool   `json:"closeTab"`
	I18nKey     string `json:"i18nKey"`
	Order       int    `json:"order"`
}
type UserMenuItemOrder []*UserMenuItem

func (m UserMenuItemOrder) Len() int           { return len(m) }
func (m UserMenuItemOrder) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m UserMenuItemOrder) Less(i, j int) bool { return m[i].Meta.Order <= m[j].Meta.Order }

// BuildMenuTree builds a tree structure from a flat list of menu items
func BuildMenuTree(menuItems []*UserMenuItem) []*UserMenuItem {
	menuMap := make(map[int]*UserMenuItem)

	for _, item := range menuItems {
		menuMap[item.ID] = item
	}

	var roots []*UserMenuItem
	for _, item := range menuItems {
		parent, exists := menuMap[item.PId]
		if exists {
			parent.Children = append(parent.Children, item)
		} else {
			roots = append(roots, item)
		}
	}
	for _, item := range menuItems {
		if len(item.Children) > 0 {
			sort.Sort(UserMenuItemOrder(item.Children))
		}
	}
	sort.Sort(UserMenuItemOrder(roots))
	return roots
}
