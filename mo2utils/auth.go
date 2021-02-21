package mo2utils

import (
	"mo2/dto"
	"mo2/server/model"

	"github.com/gin-gonic/gin"
)

// Contains test slice contain string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// GetUserInfo get user info from gin context
func GetUserInfo(ctx *gin.Context) (uinfo dto.LoginUserInfo, ext bool) {
	info, exist := ctx.Get(UserInfoKey)
	if !exist {
		ext = exist
		return
	}
	uinfo, ext = info.(dto.LoginUserInfo)

	return
}

// IsInRole check if user is in role
func IsInRole(ctx *gin.Context, role model.Erole) (result bool) {
	uinfo, ext := GetUserInfo(ctx)
	if !ext || uinfo.Roles == nil {
		return false
	}
	return Contains(uinfo.Roles, role)
}
