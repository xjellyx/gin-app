package usecase

import (
	"context"
	"time"

	"gin-app/internal/domain"

	"go.uber.org/zap"
)

type signupUsecase struct {
	repo           domain.UserRepo
	log            *zap.Logger
	contextTimeout time.Duration
}

func NewSignupUsecase(repo domain.UserRepo, timeout time.Duration, log *zap.Logger) domain.SignupUsecase {
	return &signupUsecase{
		repo:           repo,
		log:            log,
		contextTimeout: timeout,
	}
}

// Signup 注册
func (s *signupUsecase) Signup(ctx context.Context, req *domain.SignupReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}
