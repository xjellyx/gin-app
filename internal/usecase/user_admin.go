package usecase

import (
	"context"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"time"

	"gin-app/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// UserAdminConfig 用户管理员配置
type UserAdminConfig struct {
	Repo           domain.UserRepo
	ContextTimeout time.Duration
}

// userAdminUsecase 用户管理员用例
type userAdminUsecase struct {
	cfg UserAdminConfig
}

// NewUserAdminUsecase 新建用户管理员用例
func NewUserAdminUsecase(cfg UserAdminConfig) domain.UserAdminUsecase {
	return &userAdminUsecase{
		cfg: cfg,
	}
}

// List 获取用户列表
func (u *userAdminUsecase) List(ctx context.Context, req *request.UserAdminListReq) (*response.UserAdminListResp, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	data, err := u.cfg.Repo.Find(ctx, clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		},
	})
	if err != nil {
		return nil, err
	}
	ret := &response.UserAdminListResp{}
	ret.Pagination = &response.Pagination{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
	ret.Pagination.Total, err = u.cfg.Repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		ret.List = append(ret.List, &response.UserListInfo{
			ID:     v.ID,
			Uuid:   v.Uuid,
			Status: v.Status,
			UserInfo: response.UserInfo{
				Gender:    v.Gender,
				Email:     v.Email,
				Username:  v.Username,
				Phone:     v.Phone,
				CreatedAt: v.CreatedAt,
			},
		})
	}
	return ret, nil
}

// Add 添加用户
func (u *userAdminUsecase) Add(ctx context.Context, req *request.UserAdminAddReq) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	user := &domain.User{
		Uuid:     uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Gender:   req.Gender,
		Status:   req.Status,
		Phone:    req.Phone,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Phone), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	if err = u.cfg.Repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}
