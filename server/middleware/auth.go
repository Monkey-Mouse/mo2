package middleware

import (
	"mo2/mo2utils"
	"mo2/server/controller"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlware(c *gin.Context) {
	if !strings.Contains(c.Request.URL.Path, "api") || strings.Contains(c.Request.URL.Path, "/logs") {
		c.Next()
		return
	}

	cookieStr, err := c.Cookie("jwtToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	uinfo, err := mo2utils.ParseJwt(cookieStr)
	c.Set(mo2utils.UserInfoKey, uinfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	c.Next()
}
