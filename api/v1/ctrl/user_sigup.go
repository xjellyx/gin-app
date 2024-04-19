package ctrl

import (
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserSignupCtrl ctrl
type UserSignupCtrl struct {
	Usecase domain.SignupUsecase
}

func NewSignupCtl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	repo := repository.NewUserRepo(app.Database)
	us := usecase.NewSignupUsecase(usecase.SignupUsecaseConfig{
		Repo:           repo,
		ContextTimeout: timeout,
	})
	sc := UserSignupCtrl{
		Usecase: us,
	}
	group.POST("/signup", sc.Signup)
	group.POST("/sing-in", sc.SingIn)
}

// Signup
// @Tags UserSignup
// @Summary 用户注册
// @Version 1.0
// @Produce application/json
// @Param {} body domain.SignupReq true "body"
// @Router /api/v1/signup [post]
// @Success 200 {object} domain.Response
// @Security ApiKeyAuth
func (u *UserSignupCtrl) Signup(c *gin.Context) {
	var req domain.SignupReq
	var err error
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = u.Usecase.Signup(c.Request.Context(), &req); err != nil {
		return
	}

	SuccessResponse(c, nil)
}

// SingIn 登入
// @Tags UserSingIn
// @Summary 用户登入
// @Version 1.0
// @Produce application/json
// @Param {} body domain.SignInReq true "body"
// @Router /api/v1/sing-in [post]
// @Success 200 {object} domain.Response
// @Security ApiKeyAuth
func (u *UserSignupCtrl) SingIn(c *gin.Context) {
	var req domain.SignInReq
	var err error
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	res, err := u.Usecase.SignIn(c.Request.Context(), &req)
	if err != nil {
		return
	}
	SuccessResponse(c, res)
}
