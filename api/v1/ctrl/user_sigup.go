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
// @tags
// @Summary
// @Description
// @Param
// @router []
// @Success 200 {object} response.Response{}
// @Security ApiKeyAuth
// @Failure 404 {object} string
// @Failure 500 {object} string
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
