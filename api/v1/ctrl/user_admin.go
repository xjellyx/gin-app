package ctrl

import (
	"gin-app/internal/domain/request"
	"strconv"
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserAdminCtrl 管理员管理用户控制器
type UserAdminCtrl struct {
	Usecase domain.UserAdminUsecase
}

func SetupUserAdminRoute(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &UserAdminCtrl{}
	repo := repository.NewUserRepo(app.Database)
	roleRepo := repository.NewRoleRepo(app.Database)
	ctl.Usecase = usecase.NewUserAdminUsecase(usecase.UserAdminConfig{Repo: repo,
		RoleRepo:       roleRepo,
		ContextTimeout: timeout})
	h := group.Group("/users")
	h.GET("", ctl.GetUserList)
	h.POST("", ctl.AddUser)
	h.PUT("/:id", ctl.UpdateUser)
	h.DELETE("/:id", ctl.DeleteUser)
	h.DELETE("", ctl.DeleteBatchUser)
}

// GetUserList
// @Tags User 用户管理
// @Summary 用户列表
// @Version 1.0
// @Param req query request.UserAdminListReq true "查询参数"
// @Produce application/json
// @Router /api/v1/users [get]
// @Success 200 {object} response.Response{data=response.UserAdminListResp}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) GetUserList(c *gin.Context) {
	var (
		err error
		req request.UserAdminListReq
	)
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()
	if err = c.ShouldBindQuery(&req); err != nil {
		return
	}
	detail, err := u.Usecase.List(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, detail)
}

// AddUser adds a user
// @Tags User 用户管理
// @Summary 添加用户
// @Version 1.0
// @Produce application/json
// @Param req body request.UserAdminAddReq true "添加用户"
// @Router /api/v1/users [post]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) AddUser(c *gin.Context) {
	var (
		err error
		req request.UserAdminAddReq
	)
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	if err = u.Usecase.Add(c.Request.Context(), &req); err != nil {
		return
	}
	SuccessResponse(c, nil)
}

// UpdateUser updates a user
// @Tags User 用户管理
// @Summary 更新用户
// @Version 1.0
// @Produce application/json
// @Param id path int true "用户ID"
// @Param req body request.UserAdminUpdateReq true "更新用户"
// @Router /api/v1/users/{id} [put]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) UpdateUser(c *gin.Context) {
	var (
		err error
		req request.UserAdminUpdateReq
	)
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return
	}
	if err = u.Usecase.Update(c.Request.Context(), uint(id), &req); err != nil {
		return
	}
	SuccessResponse(c, nil)
}

// DeleteUser deletes a user
// @Tags User 用户管理
// @Summary 删除用户
// @Version 1.0
// @Produce application/json
// @Param id path int true "用户ID"
// @Router /api/v1/users/{id} [delete]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) DeleteUser(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return
	}
	if err = u.Usecase.Delete(c.Request.Context(), uint(id)); err != nil {
		return
	}
	SuccessResponse(c, nil)
}

// DeleteBatchUser deletes users
// @Tags User 用户管理
// @Summary 批量删除用户
// @Version 1.0
// @Produce application/json
// @Param ids body []int true "用户ID列表"
// @Router /api/v1/users/batch [delete]
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) DeleteBatchUser(c *gin.Context) {
	var (
		err error
		ids []uint
	)
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()
	if err = c.ShouldBindJSON(&ids); err != nil {
		return
	}
	if err = u.Usecase.DeleteBatch(c.Request.Context(), ids); err != nil {
		return
	}
	SuccessResponse(c, nil)
}
