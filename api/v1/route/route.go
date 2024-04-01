package route

import (
	"time"

	"gin-app/api/v1/ctrl"
	"gin-app/api/v1/middleware"
	_ "gin-app/docs"
	"gin-app/internal/bootstrap"

	"github.com/gin-gonic/gin"
	swagfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func Setup(app *bootstrap.Application, timeout time.Duration, en *gin.Engine) {
	en.Use(middleware.LimitRequestRate(app.Limiter))
	publicRouter := en.Group("/api/v1")
	publicRouter.GET("/docs/*any", ginswagger.WrapHandler(swagfiles.Handler))
	publicRouter.Use(middleware.HandlerHeadersCtx(), middleware.HandlerError())
	ctrl.NewSignupCtl(app, timeout, publicRouter)
	if app.Conf.JWTEnabled {
		publicRouter.Use(middleware.HandlerAuth())
	}
	ctrl.NewUserHimSelfCtrl(app, timeout, publicRouter)
	ctrl.NewAdminCtrl(app, timeout, publicRouter)
}
