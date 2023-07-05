package route

import (
	"time"

	"gin-app/api/v1/middleware"
	_ "gin-app/docs"
	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"

	"github.com/gin-gonic/gin"
	swagfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Setup(conf *bootstrap.Conf, timeout time.Duration, database domain.Database, log *zap.Logger, gin *gin.Engine) {
	publicRouter := gin.Group("/api/v1")
	publicRouter.GET("/docs/*any", ginswagger.WrapHandler(swagfiles.Handler))
	publicRouter.Use(middleware.HandlerHeadersCtx(), middleware.HandlerError())
	NewSignupCtl(database, timeout, log, publicRouter)
}
