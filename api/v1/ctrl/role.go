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
	"strconv"
	"time"
)

type roleController struct {
	usecase       domain.RoleUsecase
	casbinUsecase domain.CasbinRuleUsecase
}

func SetupRoleRoute(app *bootstrap.Application, timeout time.Duration, r *gin.RouterGroup) {
	ctl := &roleController{}
	repo := repository.NewRoleRepo(app.Database)
	sysApiRepo := repository.NewSysAPIRepo(app.Database)
	menuRepo := repository.NewMenuRepo(app.Database)
	ctl.usecase = usecase.NewRoleUsecase(usecase.RoleUsecaseConfig{
		SysAPIRepo: sysApiRepo,
		Casbin:     app.Casbin,
		MenuRepo:   menuRepo,
		Repo:       repo,
		Timeout:    timeout})
	ctl.casbinUsecase = usecase.NewCasbinRuleUsecase(usecase.CasbinRuleUsecaseConfig{
		SysApiRepo: sysApiRepo,
		MenuRepo:   menuRepo,
		Casbin:     app.Casbin,
		Timeout:    timeout,
	})
	h := r.Group("/roles")
	h.GET("", ctl.getRoles)
	h.GET("/all", ctl.getAllRoles)
	h.POST("", ctl.addRole)
	h.DELETE("", ctl.deleteBatch)
	h.DELETE("/:id", ctl.delete)
	h.PUT("/:id", ctl.edit)
	h.POST("perm", ctl.addApiPerm)
	h.GET(":id/perm", ctl.getApiPerm)
	h.POST("menu", ctl.addRoleMenuPerm)
	h.GET(":id/menu", ctl.getRoleMenuPerm)
	h.POST("front-page", ctl.setFrontPage)
	h.GET(":id/front-page", ctl.getFrontPage)
}

// @Tags Role 角色管理
// @Summary 获取角色列表
// @Version 1.0
// @Produce application/json
// @Param {} query request.GetRolesReq true "查询参数"
// @Router /api/v1/roles [get]
// @Success 200 {object} response.Response{data=response.GetRolesResp}
// @Security ApiKeyAuth
func (u *roleController) getRoles(c *gin.Context) {
	req := &request.GetRolesReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		_ = c.Error(err)
		return
	}
	resp, err := u.usecase.GetRoles(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, resp)
}

// @Tags Role 角色管理
// @Summary 获取全部角色
// @Version 1.0
// @Produce application/json
// @Router /api/v1/roles/all [get]
// @Success 200 {object} response.Response{data=[]response.Role}
// @Security ApiKeyAuth
func (u *roleController) getAllRoles(c *gin.Context) {
	roles, err := u.usecase.GetAllRoles(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, roles)
}

// @Tags Role 角色管理
// @Summary 添加角色
// @Version 1.0
// @Produce application/json
// @Param {} body request.AddRoleReq true "请求参数"
// @Router /api/v1/roles [post]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *roleController) addRole(c *gin.Context) {
	req := &request.AddRoleReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.AddRole(c.Request.Context(), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 批量删除角色
// @Version 1.0
// @Produce application/json
// @Param ids body []int true "请求参数"
// @Router /api/v1/roles/batch [delete]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *roleController) deleteBatch(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.DeleteBatch(c.Request.Context(), ids); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 删除角色
// @Version 1.0
// @Produce application/json
// @Param id path int true "请求参数"
// @Router /api/v1/roles/{id} [delete]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *roleController) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err = u.usecase.Delete(c.Request.Context(), uint(id)); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 编辑角色
// @Version 1.0
// @Produce application/json
// @Param id path int true "请求参数"
// @Param {} body request.EditRoleReq true "请求参数"
// @Router /api/v1/roles/{id} [put]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *roleController) edit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}
	req := &request.EditRoleReq{}
	if err = c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err = u.usecase.EditRole(c.Request.Context(), uint(id), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 配置角色api权限
// @Version 1.0
// @Produce application/json
// @Param {} body request.AddRoleAPIPermissionReq true "配置菜单api权限"
// @Router /api/v1/roles/perm [post]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) addApiPerm(c *gin.Context) {
	req := &request.AddRoleAPIPermissionReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.casbinUsecase.AddAPIPermission(c.Request.Context(), types.CasbinRoleKey, &request.AddAPIPermissionReq{
		ID:     req.RoleId,
		APIIds: req.APIIds,
	}); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 获取角色api权限
// @Version 1.0
// @Produce application/json
// @Router /api/v1/roles/{id}/perm [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) getApiPerm(c *gin.Context) {
	idStr := c.Param("id")
	data := u.casbinUsecase.GetAPIPermission(c.Request.Context(), types.CasbinRoleKey, cast.ToUint(idStr))
	SuccessResponse(c, data)
}

// @Tags Role 角色管理
// @Summary 获取角色菜单
// @Version 1.0
// @Produce application/json
// @Router /api/v1/roles/{id}/menu/tree [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) getMenu(c *gin.Context) {
	//idStr := c.Param("id")
	//data := u.casbinUsecase.GetMenu(c.Request.Context(), types.CasbinRoleKey, cast.ToUint(idStr))
	//SuccessResponse(c, data)
}

// @Tags Role 角色管理
// @Summary 配置角色菜单权限
// @Version 1.0
// @Produce application/json
// @Param {} body request.AddRoleMenuPermissionReq true "配置角色菜单权限"
// @Router /api/v1/roles/menu [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) addRoleMenuPerm(c *gin.Context) {
	req := &request.AddRoleMenuPermissionReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.casbinUsecase.AddRoleMenuPermission(c.Request.Context(), types.CasbinRoleMenuKey, req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 获取角色菜单权限
// @Version 1.0
// @Produce application/json
// @Router /api/v1/roles/{id}/menu [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) getRoleMenuPerm(c *gin.Context) {
	idStr := c.Param("id")
	data := u.casbinUsecase.GetRoleMenuPermission(c.Request.Context(), types.CasbinRoleMenuKey, cast.ToUint(idStr))
	SuccessResponse(c, data)
}

// @Tags Role 角色管理
// @Summary 设置角色前台页面
// @Version 1.0
// @Produce application/json
// @Param {} body request.SetRoleRouteFrontPageReq true "设置角色前台页面"
// @Router /api/v1/roles/front-page [post]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) setFrontPage(c *gin.Context) {
	req := &request.SetRoleRouteFrontPageReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.casbinUsecase.SetRoleRouteFrontPage(c.Request.Context(), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags Role 角色管理
// @Summary 获取角色前台页面
// @Version 1.0
// @Produce application/json
// @Router /api/v1/roles/{id}/front-page [get]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *roleController) getFrontPage(c *gin.Context) {
	idStr := c.Param("id")
	data := u.casbinUsecase.GetRoleRouteFrontPage(c.Request.Context(), cast.ToUint(idStr))
	SuccessResponse(c, data)
}
