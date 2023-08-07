package ctrl

import (
	"gin-app/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserHimSelfCtrl struct {
	Usecase domain.UserHimSelfUsecase
}

// GetUserInfo
// @Tags UserHimSelf
// @Summary 用户信息
// @Version 1.0
// @Produce application/json
// @Router /api/v1/user/info [get]
// @Success 200 {} {}
// @Security ApiKeyAuth
func (u *UserHimSelfCtrl) GetUserInfo(c *gin.Context) {

	detail, err := u.Usecase.Info(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	SuccessResponse(c, detail)
}
