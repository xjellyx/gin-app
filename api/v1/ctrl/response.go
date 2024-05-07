package ctrl

import (
	"gin-app/internal/domain/response"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功返回
func SuccessResponse(c *gin.Context, data any) {
	c.JSON(200, response.Response{Code: "SUCCESS", Msg: "操作成功", Data: data})
}
