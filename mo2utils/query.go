package mo2utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParsePagination get pagination info from ctx
func ParsePagination(ctx *gin.Context) (page int64, pagesize int64, err error) {
	page, pagesize = 0, 10
	p, ext := ctx.GetQuery("page")
	if !ext {
		return
	}
	page, err = strconv.ParseInt(p, 10, 64)
	if err != nil {
		return
	}
	ps, ext := ctx.GetQuery("pagesize")
	if !ext {
		ps, ext = ctx.GetQuery("pageSize")
		if !ext {
			return
		}
	}
	pagesize, err = strconv.ParseInt(ps, 10, 64)
	if err != nil {
		return
	}
	return
}
