package middleware

import "github.com/gin-gonic/gin"

//RoleHolder user interface with role
type RoleHolder interface {
	IsInRole(role string) bool
}

// FromCTX 从context中获取符合RoleHolder类型的用户信息的方法，如果无法获取则返回err
type FromCTX func(ctx *gin.Context) (uinfo RoleHolder, err error)

type handlerProp struct {
	handler   gin.HandlerFunc
	needRoles [][]string
	limit     int
}
type handlerKey struct {
	url    string
	method string
}
type handlerMap struct {
	innerMap   map[handlerKey]handlerProp
	prefixPath string
	roles      [][]string
	limit      int
}
