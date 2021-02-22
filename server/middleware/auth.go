package middleware

import (
	"mo2/mo2utils"
	"mo2/server/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddlware(c *gin.Context) {
	key := HandlerKey{c.FullPath(), c.Request.Method}
	prop, ok := Handlers[key]
	if !ok || prop.NeedRoles == nil || len(prop.NeedRoles) == 0 {
		c.Next()
		return
	}
	cookieStr, err := c.Cookie("jwtToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	uinfo, err := mo2utils.ParseJwt(cookieStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	for _, v := range prop.NeedRoles {
		if !mo2utils.Contains(uinfo.Roles, v) {
			c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Need role: "+v))
			return
		}
	}
	c.Set(mo2utils.UserInfoKey, uinfo)
	c.Next()
}

type HandlerProp struct {
	Handler   gin.HandlerFunc
	NeedRoles []string
}
type HandlerKey struct {
	Url    string
	Method string
}

var Handlers = make(map[HandlerKey]HandlerProp, 0)
