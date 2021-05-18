package badresponse

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	BadRequestReason  = "非法输入"
	UnauthorizeReason = "请先登录"
	NoAccessReason    = "没有权限"
)

// ResponseError 返回的错误
type ResponseError struct {
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}

// SetResponseError 通过错误获取返回的json结构
func SetResponseError(err error) (r ResponseError) {
	r.Reason = err.Error()
	r.Time = time.Now()
	return
}

// SetResponseReason 直接设置返回错误的信息
func SetResponseReason(err string) (r ResponseError) {
	r.Reason = err
	r.Time = time.Now()
	return
}

// SetErrResponse set the standard err response for aborted request.
// 记得调用之后使用return！
func SetErrResponse(ctx *gin.Context, httpStatus int, err string) {
	var r ResponseError
	r.Reason = err
	r.Time = time.Now()
	ctx.AbortWithStatusJSON(httpStatus, r)
}
