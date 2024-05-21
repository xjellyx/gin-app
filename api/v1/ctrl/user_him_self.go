package ctrl

import (
	"gin-app/api/v1/middleware"
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

func SetupUserHimSelfRoute(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &UserHimSelfCtrl{}
	repo := repository.NewUserRepo(app.Database)
	menuRepo := repository.NewMenuRepo(app.Database)
	roleRepo := repository.NewRoleRepo(app.Database)
	ctl.Usecase = usecase.NewUserHimSelfUsecase(usecase.UserHimSelfUsecaseConfig{
		Repo:           repo,
		RoleRepo:       roleRepo,
		MenuRepo:       menuRepo,
		Casbin:         app.Casbin,
		ContextTimeout: timeout})
	h := group.Group("/user")
	h.Use(middleware.HandlerAuth(true))
	h.GET("/info", ctl.GetUserInfo)
	h.GET("/menus", ctl.GetMenusTree)
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

// GetMenusTree
// @Tags UserHimSelf
// @Summary 用户菜单
// @Version 1.0
// @Produce application/json
// @Router /api/v1/user/menus [get]
// @Param code query string true "角色代码"
// @Success 200 {object} response.Response{data=response.GetMenusTreeResp}
// @Security ApiKeyAuth
func (u *UserHimSelfCtrl) GetMenusTree(c *gin.Context) {
	code := c.Query("code")
	detail, err := u.Usecase.GetMenusTree(c.Request.Context(), code)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, detail)
}
