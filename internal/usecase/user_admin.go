package usecase

import (
	"context"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/response"
	"gin-app/internal/domain/types"
	"gin-app/pkg/str"
	"time"

	"gin-app/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// UserAdminConfig 用户管理员配置
type UserAdminConfig struct {
	Repo           domain.UserRepo
	RoleRepo       domain.RoleRepo
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
	var cond []clause.Expression
	if req.Username != "" {
		cond = append(cond, clause.Eq{
			Column: "username",
			Value:  req.Username,
		})
	}
	if req.Email != "" {
		cond = append(cond, clause.Eq{
			Column: "email",
			Value:  req.Email,
		})
	}
	if req.Phone != "" {
		cond = append(cond, clause.Eq{
			Column: "phone",
			Value:  req.Phone,
		})
	}
	if req.Gender != 0 {
		cond = append(cond, clause.Eq{
			Column: "gender",
			Value:  req.Gender,
		})
	}
	if req.Status != 0 {
		cond = append(cond, clause.Eq{
			Column: "status",
			Value:  req.Status,
		})
	}
	total, err := u.cfg.Repo.Count(ctx, cond...)
	if err != nil {
		return nil, err
	}
	cond = append(cond, req.SetOrmExpression())
	cond = append(cond, clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		},
	})
	var data []*domain.User
	err = u.cfg.Repo.Database().DB(ctx).Preload("Roles").Clauses(cond...).Find(&data).Error
	if err != nil {
		return nil, err
	}
	ret := &response.UserAdminListResp{Records: make([]*response.UserListInfo, 0)}
	ret.Pagination = &response.Pagination{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
	ret.Pagination.Total = total
	for _, v := range data {
		d := &response.UserListInfo{
			ID:     v.ID,
			Uuid:   v.Uuid,
			Status: v.Status,
			UserInfo: response.UserInfo{
				Gender:    v.Gender,
				Email:     v.Email,
				Username:  v.Username,
				Phone:     v.Phone,
				CreatedAt: v.CreatedAt,
				Roles:     make([]string, 0),
			},
		}
		for _, _v := range v.Roles {
			d.Roles = append(d.Roles, _v.Code)
		}
		ret.Records = append(ret.Records, d)
	}
	return ret, nil
}

// Add 添加用户
func (u *userAdminUsecase) Add(ctx context.Context, req *request.UserAdminAddReq) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	if req.Username == "" {
		req.Username = str.GenerateRandomString(8)
	}
	if req.Gender == 0 {
		req.Gender = types.GenderUnknown
	}
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

	if len(req.Roles) > 0 {
		var roles []*domain.Role
		var values []any
		for _, v := range req.Roles {
			values = append(values, v)
		}
		roles, err = u.cfg.RoleRepo.Find(ctx, clause.IN{
			Column: "code",
			Values: values,
		})
		if err != nil {
			return err
		}
		user.Roles = roles
	}
	if err = u.cfg.Repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// Update 更新用户
func (u *userAdminUsecase) Update(ctx context.Context, id uint, req *request.UserAdminUpdateReq) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()

	err := u.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		user, err := u.cfg.Repo.FindOne(ctx, id)
		if err != nil {
			return err
		}
		err = u.cfg.Repo.Update(ctx, id, &domain.User{
			Username: req.Username,
			Email:    req.Email,
			Gender:   req.Gender,
			Status:   req.Status,
			Phone:    req.Phone,
		})
		if err != nil {
			return err
		}
		if err = u.cfg.Repo.Database().DB(ctx).Model(user).Association("Roles").Clear(); err != nil {
			return err
		}
		if len(req.Roles) > 0 {
			var roles []*domain.Role
			var values []any
			for _, v := range req.Roles {
				values = append(values, v)
			}
			roles, err = u.cfg.RoleRepo.Find(ctx, clause.IN{
				Column: "code",
				Values: values,
			})
			if err != nil {
				return err
			}
			if err = u.cfg.Repo.Database().DB(ctx).Model(user).Association("Roles").Append(roles); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Delete 删除用户
func (u *userAdminUsecase) Delete(ctx context.Context, id uint) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	user, err := u.cfg.Repo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	// 删除关联的角色信息
	if err = u.cfg.Repo.Database().DB(ctx).Model(user).Association("Roles").Clear(); err != nil {
		return err
	}
	return u.cfg.Repo.DeleteOne(ctx, id)
}

func (u *userAdminUsecase) DeleteBatch(ctx context.Context, ids []uint) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	cond := clause.IN{}
	for _, v := range ids {
		cond.Values = append(cond.Values, v)
	}
	cond.Column = "id"
	err := u.cfg.Repo.ExecTX(ctx, func(ctx context.Context) error {
		err := u.cfg.Repo.Database().DB(ctx).Table("user_roles").Clauses(clause.IN{Column: "user_id", Values: cond.Values}).Delete(ctx).Error
		if err != nil {
			return err
		}
		return u.cfg.Repo.DeleteBy(ctx, []clause.Expression{cond})
	})
	if err != nil {
		return err
	}
	return nil
}
