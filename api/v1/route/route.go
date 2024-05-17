package route

import (
	"fmt"
	"net/http"
	"time"

	"gin-app/api/v1/ctrl"
	"gin-app/api/v1/middleware"
	_ "gin-app/docs"
	"gin-app/internal/bootstrap"

	"github.com/gin-gonic/gin"
	swagfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func Setup(app *bootstrap.Application, timeout time.Duration) *http.Server {
	en := gin.Default()
	en.Use(middleware.LimitRequestRate(app.Limiter))
	publicRouter := en.Group("/api/v1")
	publicRouter.GET("/docs/*any", ginswagger.WrapHandler(swagfiles.Handler))
	publicRouter.Use(middleware.HandlerHeadersCtx(), middleware.HandlerError())
	ctrl.SetupUserSignRoute(app, timeout, publicRouter)
	if app.Conf.JWT.Enable {
		publicRouter.Use(middleware.HandlerAuth())
	}
	ctrl.SetupUserHimSelfRoute(app, timeout, publicRouter)
	ctrl.SetupUserAdminRoute(app, timeout, publicRouter)
	ctrl.SetupRoleRoute(app, timeout, publicRouter)
	ctrl.SetupMenuRoute(app, timeout, publicRouter)
	ctrl.SetupSysAPIRoute(app, timeout, publicRouter)
	srv := &http.Server{
		Addr:    fmt.Sprintf(`:%v`, app.Conf.HTTPort),
		Handler: en,
	}
	return srv
}
