package route

import (
	"gin-app/api/v1/ctrl"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewSignupCtl(data domain.Database, timeout time.Duration, log *zap.Logger, group *gin.RouterGroup) {
	repo := repository.NewUserRepo(data)
	us := usecase.NewSignupUsecase(repo, timeout, log)
	sc := ctrl.UserSignupCtrl{
		Usecase: us,
	}
	group.POST("/signup", sc.Signup)
}
