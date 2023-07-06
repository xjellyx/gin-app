package usecase

import (
	"context"
	"time"

	"gin-app/internal/domain"

	"go.uber.org/zap"
)

type SignupUsecaseConfig struct {
	Repo           domain.UserRepo
	Log            *zap.Logger
	ContextTimeout time.Duration
}

type signupUsecase struct {
	cfg SignupUsecaseConfig
}

func NewSignupUsecase(cfg SignupUsecaseConfig) domain.SignupUsecase {
	return &signupUsecase{
		cfg: cfg,
	}
}

// Signup 注册
func (s *signupUsecase) Signup(ctx context.Context, req *domain.SignupReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ContextTimeout)
	defer cancel()
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.cfg.Repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}
