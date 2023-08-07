package route

import (
	"time"

	"gin-app/api/v1/ctrl"
	"gin-app/internal/bootstrap"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewAdminCtrl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &ctrl.UserAdminCtrl{}
	repo := repository.NewUserRepo(app.Database)
	ctl.Usecase = usecase.NewUserAdminUsecase(usecase.UserAdminConfig{Repo: repo, ContextTimeout: timeout})
	h := group.Group("/users")
	h.GET("", ctl.GetUserList)
	h.POST("", ctl.AddUser)
}
