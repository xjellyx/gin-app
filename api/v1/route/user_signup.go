package route

import (
	"time"

	"gin-app/api/v1/ctrl"
	"gin-app/internal/bootstrap"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewSignupCtl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	repo := repository.NewUserRepo(app.Database)
	us := usecase.NewSignupUsecase(usecase.SignupUsecaseConfig{
		Repo:           repo,
		Cache:          app.Rdb,
		ContextTimeout: timeout,
	})
	sc := ctrl.UserSignupCtrl{
		Usecase: us,
	}
	group.POST("/signup", sc.Signup)
	group.POST("/sing-in", sc.SingIn)
}
