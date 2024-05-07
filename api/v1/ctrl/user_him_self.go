package ctrl

import (
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHimSelfCtrl struct {
	Usecase domain.UserHimSelfUsecase
}

func NewUserHimSelfCtrl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &UserHimSelfCtrl{}
	repo := repository.NewUserRepo(app.Database)
	ctl.Usecase = usecase.NewUserHimSelfUsecase(usecase.UserHimSelfUsecaseConfig{Repo: repo, ContextTimeout: timeout})
	h := group.Group("/user")
	h.GET("/info", ctl.GetUserInfo)
}

// GetUserInfo
// @Tags UserHimSelf
// @Summary 用户信息
// @Version 1.0
// @Produce application/json
// @Router /api/v1/user/info [get]
// @Success 200 {object} response.Response{data=response.UserInfo}
// @Security ApiKeyAuth
func (u *UserHimSelfCtrl) GetUserInfo(c *gin.Context) {

	detail, err := u.Usecase.Info(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, detail)
}
