package usecase

import (
	"context"
	"time"

	"gin-app/internal/domain"
	"gin-app/pkg/scontext"

	"gorm.io/gorm/clause"
)

type UserHimSelfUsecaseConfig struct {
	Repo           domain.UserRepo
	ContextTimeout time.Duration
}

type userHimSelfUsecase struct {
	cfg UserHimSelfUsecaseConfig
}

func NewUserHimSelfUsecase(cfg UserHimSelfUsecaseConfig) domain.UserHimSelfUsecase {
	return &userHimSelfUsecase{cfg: cfg}
}

func (u *userHimSelfUsecase) Info(ctx context.Context) (*domain.UserInfo, error) {
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
	data := &domain.UserInfo{
		Username:  ret.Username,
		Email:     ret.Email,
		Phone:     ret.Phone,
		Gender:    ret.Gender,
		CreatedAt: ret.CreatedAt,
	}
	return data, nil
}
