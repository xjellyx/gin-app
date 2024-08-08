package ctrl

import (
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/domain/request"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type departmentController struct {
	usecase domain.DepartmentUsecase
}

func SetupDepartmentRoute(app *bootstrap.Application, timeout time.Duration, r *gin.RouterGroup) {
	ctl := &departmentController{}
	repo := repository.NewDepartmentRepo(app.Database)
	ctl.usecase = usecase.NewDepartmentUsecase(&usecase.DepartmentUsecaseCfg{DepartmentRepo: repo})
	h := r.Group("/departments")
	h.POST("", ctl.addDepartment)
	h.GET("", ctl.getDepartments)
	h.DELETE(":id", ctl.deleteDepartment)
	h.PUT(":id", ctl.editDepartment)
	h.GET("tree", ctl.getDepartmentTree)
}

// @Tags department 部门管理
// @Summary 新增部门
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Param {} body request.AddDepartmentReq true "新增部门"
// @Router /api/v1/departments [post]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *departmentController) addDepartment(c *gin.Context) {
	req := &request.AddDepartmentReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.AddDepartment(c.Request.Context(), req); err != nil {
		_ = c.Error(err)
	}
	SuccessResponse(c, nil)
}

// @Tags department 部门管理
// @Summary 获取部门列表
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Param {} query request.GetDepartmentsReq true "获取部门列表"
// @Param pageSize query int false "页大小"
// @Param pageNum query int false "页码"
// @Router /api/v1/departments [get]
// @Success 200 {object} response.Response{data=response.GetDepartmentsResp}
// @Security ApiKeyAuth
func (u *departmentController) getDepartments(c *gin.Context) {
	req := &request.GetDepartmentsReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		_ = c.Error(err)
		return
	}
	data, err := u.usecase.GetDepartments(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, data)
}

// @Tags department 部门管理
// @Summary 删除部门
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "部门ID"
// @Router /api/v1/departments/{id} [delete]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *departmentController) deleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	if err := u.usecase.Delete(c.Request.Context(), cast.ToUint(idStr)); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags department 部门管理
// @Summary 编辑部门
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "部门ID"
// @Param {} body request.EditDepartmentReq true "编辑部门"
// @Router /api/v1/departments/{id} [put]
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
func (u *departmentController) editDepartment(c *gin.Context) {
	req := &request.EditDepartmentReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := u.usecase.EditDepartment(c.Request.Context(), cast.ToUint(c.Param("id")), req); err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, nil)
}

// @Tags department 部门管理
// @Summary 获取部门树
// @Version 1.0
// @Accept application/json
// @Produce application/json
// @Router /api/v1/departments/tree [get]
// @Success 200 {object} response.Response{data=response.GetDepartmentTreeResp}
// @Security ApiKeyAuth
func (u *departmentController) getDepartmentTree(c *gin.Context) {
	data, err := u.usecase.GetDepartmentTree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, data)
}
