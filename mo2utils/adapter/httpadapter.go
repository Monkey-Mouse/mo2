package adapter

import (
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/gin-gonic/gin"
)

// RequestHandler used with ResponseAdapter
type RequestHandler func(ctx *gin.Context) (status int, json interface{}, err error)

// ResponseAdapter for gin handler
func ResponseAdapter(handler RequestHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, j, err := handler(ctx)
		processResult(ctx, s, j, err)
	}
}

func processResult(ctx *gin.Context, status int, json interface{}, err error) {
	if err != nil {
		ctx.AbortWithStatusJSON(status, badresponse.SetResponseError(err))
		return
	}
	ctx.JSON(status, json)
}
