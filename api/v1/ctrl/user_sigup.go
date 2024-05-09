package ctrl

import (
	"gin-app/internal/domain/request"
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserSignCtrl ctrl
type UserSignCtrl struct {
	Usecase domain.SignupUsecase
}

func NewUserSignCtl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	repo := repository.NewUserRepo(app.Database)
	us := usecase.NewSignupUsecase(usecase.SignupUsecaseConfig{
		Repo:           repo,
		ContextTimeout: timeout,
	})
	sc := UserSignCtrl{
		Usecase: us,
	}
	group.POST("/signup", sc.Signup)
	group.POST("/sign-in", sc.SingIn)
}

// Signup
// @Tags UserSign 用户注册登录
// @Summary 用户注册
// @Version 1.0
// @Produce application/json
// @Param {} body request.SignupReq true "body"
// @Router /api/v1/signup [post]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *UserSignCtrl) Signup(c *gin.Context) {
	var req request.SignupReq
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
// @Tags UserSign 用户注册登录
// @Summary 用户登入
// @Version 1.0
// @Produce application/json
// @Param {} body request.SignInReq true "body"
// @Router /api/v1/sign-in [post]
// @Success 200 {object} response.Response{data=response.SignInResp}
// @Security ApiKeyAuth
func (u *UserSignCtrl) SingIn(c *gin.Context) {
	var req request.SignInReq
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
