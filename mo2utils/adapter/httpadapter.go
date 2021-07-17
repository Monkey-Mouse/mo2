package adapter

import (
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/mo2utils/basiclog"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/gin-gonic/gin"
)

// RequestHandler used with ResponseAdapter
type RequestHandler func(ctx *gin.Context) (status int, json interface{}, err error)

// UserRequestHandler used with ReAdapterWithUinfo
type UserRequestHandler func(ctx *gin.Context, uinfo dto.LoginUserInfo) (status int, json interface{}, err error)

// ReAdapter for gin handler
func ReAdapter(handler RequestHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, j, err := handler(ctx)
		processResult(ctx, s, j, err)
	}
}

// ResponseAdapter for gin handler
func ReAdapterWithUinfo(handler UserRequestHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uinfo, _ := mo2utils.GetUserInfo(ctx)
		s, j, err := handler(ctx, uinfo)
		processResult(ctx, s, j, err)
	}
}

func processResult(ctx *gin.Context, status int, json interface{}, err error) {
	if status == 0 {
		status = 200
	}
	if err != nil {
		basiclog.ErrLogger.Println(err)
		ctx.AbortWithStatusJSON(status, badresponse.SetResponseError(err))
		return
	}
	ctx.JSON(status, json)
}
