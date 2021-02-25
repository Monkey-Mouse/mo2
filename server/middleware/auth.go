package middleware

import (
	"mo2/mo2utils"
	"mo2/server/controller"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func AuthMiddlware(c *gin.Context) {
	cookieStr, err := c.Cookie("jwtToken")
	uinfo, jwterr := mo2utils.ParseJwt(cookieStr)
	c.Set(mo2utils.UserInfoKey, uinfo)
	key := HandlerKey{c.FullPath(), c.Request.Method}
	prop, ok := handlers[key]
	if !ok || prop.NeedRoles == nil || len(prop.NeedRoles) == 0 {
		c.Next()
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	if jwterr != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Unauthorized!"))
		return
	}
	for _, v := range prop.NeedRoles {
		if !mo2utils.Contains(uinfo.Roles, v) {
			c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("Need role: "+v))
			return
		}
	}
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
type HandlerMap struct {
	Map        *map[HandlerKey]HandlerProp
	PrefixPath string
}

var handlers = make(map[HandlerKey]HandlerProp, 0)
var H = HandlerMap{&handlers, ""}

func (h HandlerMap) Group(relativPath string) HandlerMap {

	h.PrefixPath = path.Join(h.PrefixPath, relativPath)
	return h
}

func (h HandlerMap) Handle(method string, relativPath string, handler gin.HandlerFunc, roles ...string) {
	(*h.Map)[HandlerKey{Url: path.Join(h.PrefixPath, relativPath), Method: method}] = HandlerProp{
		Handler: handler, NeedRoles: roles}
}
func (h HandlerMap) Get(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodGet, relativPath, handler, roles...)
}
func (h HandlerMap) Post(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPost, relativPath, handler, roles...)
}
func (h HandlerMap) Delete(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodDelete, relativPath, handler, roles...)
}
func (h HandlerMap) Put(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPut, relativPath, handler, roles...)
}
func (h HandlerMap) RegisterMapedHandlers(r *gin.Engine) {
	for k, v := range *h.Map {
		r.Handle(k.Method, k.Url, v.Handler)
	}
}
