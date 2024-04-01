package ctrl

import (
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

func NewAdminCtrl(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &UserAdminCtrl{}
	repo := repository.NewUserRepo(app.Database)
	ctl.Usecase = usecase.NewUserAdminUsecase(usecase.UserAdminConfig{Repo: repo, ContextTimeout: timeout})
	h := group.Group("/users")
	h.GET("", ctl.GetUserList)
	h.POST("", ctl.AddUser)
}

// GetUserList
// @Tags UserAdminCtrl
// @Summary 用户列表
// @Version 1.0
// @Produce application/json
// @Router /api/v1/users [get]
// @Success 200 {object} domain.Response{data=domain.UserAdminListResp}
// @Security ApiKeyAuth
func (u *UserAdminCtrl) GetUserList(c *gin.Context) {
	var (
		err error
		req domain.UserAdminListReq
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
// @Tags UserAdminCtrl
// @Summary 添加用户
// @Version 1.0
// @Produce application/json
// @Router /api/v1/users [post]
// @Success 200 {object} domain.Response
// @Security ApiKeyAuth
func (u *UserAdminCtrl) AddUser(c *gin.Context) {
	var (
		err error
		req domain.UserAdminAddReq
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
