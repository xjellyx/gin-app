package ctrl

import (
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/repository"
	"gin-app/internal/usecase"
	"github.com/gin-gonic/gin"
)

type sysAPICtrl struct {
	usecase domain.SysAPIUsecase
}

func SetupSysAPIRoute(app *bootstrap.Application, timeout time.Duration, group *gin.RouterGroup) {
	ctl := &sysAPICtrl{
		usecase: usecase.NewSysAPIUsecase(usecase.SysAPIConfig{
			Repo:    repository.NewSysAPIRepo(app.Database),
			Timeout: timeout,
		}),
	}
	h := group.Group("/sys-api")
	h.GET("tree", ctl.getSysAPITree)
}

// @Tags sys-api 系统接口
// @Summary 获取系统接口树
// @Version 1.0
// @Produce application/json
// @Router /api/v1/sys-api/tree [get]
// @Success 200 {object} response.Response{data=response.SysAPITreeResp}
// @Security ApiKeyAuth
func (u *sysAPICtrl) getSysAPITree(c *gin.Context) {
	tree, err := u.usecase.GetSysAPITree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, tree)
}
