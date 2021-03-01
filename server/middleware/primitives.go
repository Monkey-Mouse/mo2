package middleware

import "github.com/gin-gonic/gin"

//RoleHolder user interface with role
type RoleHolder interface {
	IsInRole(role string) bool
}
type FromCTX func(ctx *gin.Context) (uinfo RoleHolder, err error)
