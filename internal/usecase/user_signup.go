package usecase

import (
	"context"
	"strings"
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/infra/cache"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type SignupUsecaseConfig struct {
	Repo           domain.UserRepo
	Log            *zap.Logger
	ContextTimeout time.Duration
	Cache          cache.Cache
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
		//Username:     req.Username,
		Phone: req.Phone,
		Email: req.Email,
		Uuid:  uuid.New().String(),
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(pwd)
	if err := s.cfg.Repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// SingIn 登入
func (s *signupUsecase) SingIn(ctx context.Context, req *domain.SingInReq) (*domain.SingInResp, error) {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ContextTimeout)
	defer cancel()
	var cond clause.Eq
	switch {
	case req.Email != "":
		cond = clause.Eq{
			Column: "email",
			Value:  req.Email,
		}
	case req.Phone != "":
		cond = clause.Eq{
			Column: "phone",
			Value:  req.Phone,
		}
	}
	user, err := s.cfg.Repo.FindOneBy(ctx, []clause.Expression{cond})
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}
	lan := scontext.GetLanguage(ctx)
	// 判断用户的状态
	if user.Status != domain.StatusNormal {
		switch user.Status {
		case domain.StatusLocked:
			err = serror.Error(serror.ErrUserInactivate, lan)
			return nil, err
		case domain.StatusDeleted:
			err = serror.Error(serror.ErrUserDeleted, lan)
			return nil, err
		case domain.StatusFreeze:
			err = serror.Error(serror.ErrUserFreeze, lan)
			return nil, err
		default:
			err = serror.Error(serror.ErrUserStatusAbnormal, lan)
			return nil, err
		}
	}
	res, err := generateToken(ctx, s.cfg.Cache, &Claims{
		UserUuid: user.Uuid,
		Username: user.Username,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Claims 自定义存储一些信息
type Claims struct {
	UserUuid string `json:"userUuid,omitempty"`
	Username string `json:"username,omitempty"`
	CacheKey string `json:"cacheKey,omitempty"`
	jwt.RegisteredClaims
}

func createToken(expireTime time.Duration, key string, cla *Claims) (string, error) {
	now := time.Now()
	cfg := bootstrap.GetConfig()
	cla.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expireTime * time.Minute)),
		NotBefore: jwt.NewNumericDate(now),
		ID:        cla.UserUuid,
	}

	var (
		signingMethod jwt.SigningMethod
	)
	switch strings.ToLower(cfg.JWTSigningMethod) {
	case "hs384":
		signingMethod = jwt.SigningMethodHS384
	case "hs256":
		signingMethod = jwt.SigningMethodHS256
	default:
		signingMethod = jwt.SigningMethodHS512
	}
	token := jwt.NewWithClaims(signingMethod, cla)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateToken(ctx context.Context, c cache.Cache, cla *Claims) (ret *domain.SingInResp, err error) {
	cfg := bootstrap.GetConfig()
	cla.CacheKey = cla.UserUuid
	tokenString, err := createToken(cfg.JWTExpireTime*time.Minute, string([]byte(cfg.JWTSigningKey)), cla)
	if err != nil {
		return nil, err
	}
	if err = c.Set(ctx, cla.CacheKey, tokenString, cfg.JWTExpireTime*time.Minute); err != nil {
		return nil, err
	}
	refreshCla := *cla
	refreshCla.CacheKey = cla.UserUuid + "_refresh_token"
	refreshToken, err := createToken(cfg.JWTRefreshExpireTime*time.Minute, cfg.JWTRefreshSingingKey, &refreshCla)
	if err != nil {
		return nil, err
	}
	if err = c.Set(ctx, refreshCla.CacheKey, refreshToken, cfg.JWTRefreshExpireTime*time.Minute); err != nil {
		return nil, err
	}

	tokenInfo := &domain.SingInResp{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    cla.ExpiresAt.Time,
	}
	return tokenInfo, nil
}

func jwtTokenKeyFunc(key []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
}

// ParseToken 解析令牌
func ParseToken(tokenString string, cla *Claims, key string) bool {
	token, err := jwt.ParseWithClaims(tokenString, cla, jwtTokenKeyFunc([]byte(key)))
	if err != nil {
		if err != nil {
			bootstrap.GlobalLog.Error("parseToken", zap.Error(err))
		}
		return false
	}
	return token.Valid
}
