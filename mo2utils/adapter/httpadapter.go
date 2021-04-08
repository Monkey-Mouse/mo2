package adapter

import "github.com/gin-gonic/gin"

type RequestHandler func(ctx *gin.Context) (status int, json interface{}, err error)

func ResponseAdapter(handler RequestHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, j, err := handler(ctx)
		processResult(ctx, s, j, err)
	}
}

func processResult(ctx *gin.Context, status int, json interface{}, err error) {
	if err != nil {
		ctx.AbortWithStatusJSON(status, json)
		return
	}
	ctx.JSON(status, json)
}
