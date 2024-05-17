package ctrl

import (
	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/domain/request"
	"gin-app/internal/domain/types"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"time"
)

type menuController struct {
	usecase       domain.MenuUsecase
	casbinUsecase domain.CasbinRuleUsecase
}

func SetupMenuRoute(app *bootstrap.Application, timeout time.Duration, r *gin.RouterGroup) {
	ctl := &menuController{}
	repo := repository.NewMenuRepo(app.Database)
	sysAPIRepo := repository.NewSysAPIRepo(app.Database)
	ctl.usecase = usecase.NewMenuUsecase(usecase.MenuUsecaseConfig{Repo: repo,
		Casbin:     app.Casbin,
		SysApiRepo: sysAPIRepo,
		Timeout:    timeout})
	ctl.casbinUsecase = usecase.NewCasbinRuleUsecase(usecase.CasbinRuleUsecaseConfig{SysApiRepo: sysAPIRepo, Timeout: timeout})
	h := r.Group("/menus")
	h.POST("", ctl.addMenu)
	h.GET("", ctl.getMenus)
	h.DELETE(":id", ctl.deleteMenu)
	h.DELETE("", ctl.deleteMenus)
	h.PUT(":id", ctl.editMenu)
	h.GET("tree", ctl.getMenuTree)
	h.GET("pages", ctl.getAllPages)
	h.POST("perm", ctl.addApiPerm)
	h.GET(":id/perm", ctl.getApiPerm)
	h.GET("route/exist", ctl.checkRoute)
}

// @Tags menu 菜单管理
// @Summary 新增菜单
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Param {} body request.AddMenuReq true "新增菜单"
// @Router /api/v1/menus [post]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *menuController) addMenu(c *gin.Context) {
	req := &request.AddMenuReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.AddMenu(c.Request.Context(), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags menu 菜单管理
// @Summary 获取菜单列表
// @Version 1.0
// @Produce application/json
// @Param {} query request.GetMenusReq true "查询参数"
// @Router /api/v1/menus [get]
// @Success 200 {object} response.Response{data=response.GetMenusResp}
// @Security ApiKeyAuth
func (u *menuController) getMenus(c *gin.Context) {
	req := &request.GetMenusReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		_ = c.Error(err)
		return
	}
	resp, err := u.usecase.GetMenus(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, resp)
}

// @Tags menu 菜单管理
// @Summary 删除菜单
// @Version 1.0
// @Produce application/json
// @Param id path int true "菜单ID"
// @Router /api/v1/menus/{id} [delete]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *menuController) deleteMenu(c *gin.Context) {
	menuId := c.Param("id")
	if err := u.usecase.DeleteMenu(c.Request.Context(), cast.ToUint(menuId)); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags menu 菜单管理
// @Summary 批量删除菜单
// @Version 1.0
// @Produce application/json
// @Param ids body []int true "菜单ID列表"
// @Router /api/v1/menus [delete]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *menuController) deleteMenus(c *gin.Context) {
	var (
		ids []uint
	)
	if err := c.ShouldBindJSON(&ids); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.DeleteMenus(c.Request.Context(), ids); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags menu 菜单管理
// @Summary 编辑菜单
// @Version 1.0
// @Produce application/json
// @Param id path int true "菜单ID"
// @Param {} body request.EditMenuReq true "编辑菜单"
// @Router /api/v1/menus/{id} [put]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *menuController) editMenu(c *gin.Context) {
	menuId := c.Param("id")
	req := &request.EditMenuReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.EditMenu(c.Request.Context(), cast.ToUint(menuId), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags menu 菜单管理
// @Summary 获取菜单树
// @Version 1.0
// @Produce application/json
// @Router /api/v1/menus/tree [get]
// @Success 200 {object} response.Response{data=[]response.MenuTreeItem}
// @Security ApiKeyAuth
func (u *menuController) getMenuTree(c *gin.Context) {
	resp, err := u.usecase.GetMenuTree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, resp)
}

// @Tags menu 菜单管理
// @Summary 获取全部页面
// @Version 1.0
// @Produce application/json
// @Router /api/v1/menus/pages [get]
// @Success 200 {object} response.Response{data=[]string}
// @Security ApiKeyAuth
func (u *menuController) getAllPages(c *gin.Context) {
	pages, err := u.usecase.GetAllPages(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, pages)
}

// @Tags menu 菜单管理
// @Summary 配置菜单api权限
// @Version 1.0
// @Produce application/json
// @Param {} body request.AddMenuAPIPermissionReq true "配置菜单api权限"
// @Router /api/v1/menus/perm [post]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *menuController) addApiPerm(c *gin.Context) {
	req := &request.AddMenuAPIPermissionReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.casbinUsecase.AddAPIPermission(c.Request.Context(), types.CasbinMenuKey, &request.AddAPIPermissionReq{
		ID:     req.MenuId,
		APIIds: req.APIIds,
	}); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags menu 菜单管理
// @Summary 获取菜单api权限
// @Version 1.0
// @Produce application/json
// @Router /api/v1/menus/{id}/perm [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *menuController) getApiPerm(c *gin.Context) {
	idStr := c.Param("id")
	data := u.casbinUsecase.GetAPIPermission(c.Request.Context(), types.CasbinMenuKey, cast.ToUint(idStr))
	SuccessResponse(c, data)
}

// @Tags menu 菜单管理
// @Summary 校验菜单是否存在
// @Version 1.0
// @Produce application/json
// @Param routeName query string true "菜单路由"
// @Router /api/v1/menus/route/exist [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *menuController) checkRoute(c *gin.Context) {
	routeName := c.Query("routeName")
	f := u.usecase.CheckRouteExist(c.Request.Context(), routeName)
	SuccessResponse(c, f)
}
