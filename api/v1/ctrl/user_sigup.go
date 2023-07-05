package ctrl

import (
	"gin-app/internal/domain"

	"github.com/gin-gonic/gin"
)

// UserSignupCtrl ctrl
type UserSignupCtrl struct {
	Usecase domain.SignupUsecase
}

// Signup
// @Tags UserSignup
// @Summary 用户注册
// @Version 1.0
// @Produce application/json
// @Param {} body domain.SignupReq true "body"
// @Router /api/v1/signup [post]
// @Success 200 {} {}
// @Security ApiKeyAuth
func (u *UserSignupCtrl) Signup(c *gin.Context) {
	var req domain.SignupReq
	var err error
	defer func() {
		if err != nil {
			_ = c.Error(err)
			return
		}
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = u.Usecase.Signup(c.Request.Context(), &req); err != nil {
		return
	}

	SuccessResponse(c, gin.H{"a": "b"})
}
