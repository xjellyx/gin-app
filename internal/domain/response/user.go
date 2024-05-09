package response

import (
	"time"

	"gin-app/internal/domain/types"
)

type UserInfo struct {
	Username  string           `json:"username"`  // username
	Email     string           `json:"email"`     // 邮箱
	Phone     string           `json:"phone"`     // 电话
	Gender    types.UserGender `json:"gender"`    // 性别 1未知、2男性、3女性
	CreatedAt time.Time        `json:"createdAt"` // 注册时间
}

type UserListInfo struct {
	Uuid   string           `json:"uuid"`
	ID     uint             `json:"id"`
	Status types.UserStatus `json:"status"`
	UserInfo
}

// UserAdminListResp 用户列表响应
type UserAdminListResp struct {
	List       []*UserListInfo `json:"list"`
	Pagination *Pagination     `json:"pagination"`
}

// SignInResp 用户登录响应
type SignInResp struct {
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expiresAt"`
	RefreshToken string    `json:"refreshToken"`
}
