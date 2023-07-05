package route

import (
	"time"

	"gin-app/api/v1/middleware"
	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(conf *bootstrap.Conf, timeout time.Duration, database domain.Database, log *zap.Logger, gin *gin.Engine) {
	publicRouter := gin.Group("")
	publicRouter.Use(middleware.HandlerHeadersCtx(), middleware.HandlerError())
	NewSignupCtl(database, timeout, log, publicRouter)
}
