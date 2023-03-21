package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// status code
const (
	//Success      = 200
	Success      = 0
	Failed       = 500
	TokenExpired = 401
)

// errors
const (
	ErrServerInternal = "服务器内部错误"
	ErrUserNotLogin   = "用户未登录"
	ErrUserNotExists  = "用户不存在"
	ErrUserExists     = "用户已存在"
	ErrUserPassword   = "密码错误"
	ErrUserToken      = "用户token错误"
	ErrUserTokenExp   = "用户token过期"
	ErrInvalidParams  = "参数错误"
)

// success response
const (
	SuccessMsg      = "success"
	QuerySuccessMsg = "查询成功"
	EmptyVideoList  = "暂无视频"
)

type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// 通用响应错误
func ResponseError(c *gin.Context, errMsg string) {
	var res CommonResponse
	res.StatusCode = Failed
	res.StatusMsg = errMsg
	c.JSON(http.StatusOK, res)
}

// 通用响应成功
func ResponseSuccess(c *gin.Context, msg string) {
	var res CommonResponse
	res.StatusCode = Success
	res.StatusMsg = msg
	c.JSON(http.StatusOK, res)
}
