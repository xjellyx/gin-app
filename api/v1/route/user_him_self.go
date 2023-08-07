package route

import (
	"time"

	"gin-app/api/v1/ctrl"
	"gin-app/internal/bootstrap"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewUserHimSelfCtrl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &ctrl.UserHimSelfCtrl{}
	repo := repository.NewUserRepo(app.Database)
	ctl.Usecase = usecase.NewUserHimSelfUsecase(usecase.UserHimSelfUsecaseConfig{Repo: repo, ContextTimeout: timeout})
	h := group.Group("/user")
	h.GET("/info", ctl.GetUserInfo)
}
