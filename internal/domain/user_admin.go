package domain

import "context"

// UserAdminListReq 用户列表请求
type UserAdminListReq struct {
	QueryReq
}

type UserListInfo struct {
	Uuid   string     `json:"uuid"`
	ID     uint       `json:"id"`
	Status UserStatus `json:"status"`
	UserInfo
}

// UserAdminListResp 用户列表响应
type UserAdminListResp struct {
	List       []*UserListInfo `json:"list"`
	Pagination *Pagination     `json:"pagination"`
}

type UserAdminAddReq struct {
	Username string     `json:"username" binding:"required"`           // 用户名
	Email    string     `json:"email" binding:"required,email"`        // 邮箱
	Phone    string     `json:"phone" binding:"required"`              // 手机号码
	Gender   UserGender `json:"gender" binding:"required,min=1,max=3"` // 性别1未知，2男，3女
	Status   UserStatus `json:"status" binding:"required,min=1,max=4"` // 状态：1正常，2锁定，3冻结，4删除
}

type UserAdminUsecase interface {
	List(ctx context.Context, req *UserAdminListReq) (*UserAdminListResp, error)
	Add(ctx context.Context, req *UserAdminAddReq) error
}
