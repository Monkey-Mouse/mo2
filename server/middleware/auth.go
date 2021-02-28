package middleware

import (
	"mo2/mo2utils"
	"mo2/server/controller"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modern-go/concurrent"
	"github.com/willf/bloom"
)

var duration int = 10
var unblockEvery int = 3600

// SetupRateLimiter setup ddos banner
func SetupRateLimiter(limitEvery int, unblockevery int) {
	duration = limitEvery
	unblockEvery = unblockevery
}
func cleaner() {
	for {
		for _, v := range handlers {

			v.rates = concurrent.NewMap()
		}
		time.Sleep(time.Second * time.Duration(duration))
	}
}
func resetBlocker() {
	for {
		blockFilter.ClearAll()
		time.Sleep(time.Second * time.Duration(unblockEvery))
	}
}

var blockFilter = bloom.NewWithEstimates(100000, 0.01)

func checkRateLimit(prop handlerProp, ip string) bool {
	// rate limit logic
	if prop.limit < 0 {
		return true
	}
	v, ext := prop.rates.Load(ip)
	if !ext {
		prop.rates.Store(ip, 1)
	} else {
		prop.rates.Store(ip, (v.(int))+1)
		if v.(int)+1 > prop.limit {
			blockFilter.AddString(ip)
			return false
		}
	}
	return true
}

// AuthMiddleware also have rate limit function
func AuthMiddleware(c *gin.Context) {
	// Block illegal ips
	if blockFilter.TestString(c.ClientIP()) {
		c.AbortWithStatusJSON(http.StatusForbidden, controller.SetResponseReason("IP Blocked!检测到该ip地址存在潜在的ddos行为"))
		return
	}

	cookieStr, err := c.Cookie("jwtToken")
	uinfo, jwterr := mo2utils.ParseJwt(cookieStr)
	c.Set(mo2utils.UserInfoKey, uinfo)
	key := handlerKey{c.FullPath(), c.Request.Method}
	prop, ok := handlers[key]
	// not registered for this middleware
	if !ok {
		c.Next()
		return
	}
	// rate limit logic
	if !checkRateLimit(prop, c.ClientIP()) {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, controller.SetResponseReason("Too frequent!"))
		return
	}
	// role auth logic
	if prop.NeedRoles == nil || len(prop.NeedRoles) == 0 {
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

type handlerProp struct {
	Handler   gin.HandlerFunc
	NeedRoles []string
	rates     *concurrent.Map
	limit     int
}
type handlerKey struct {
	URL    string
	Method string
}
type handlerMap struct {
	Map        map[handlerKey]handlerProp
	PrefixPath string
}

var handlers = make(map[handlerKey]handlerProp, 0)

// H handlermap, like gin router
var H = handlerMap{handlers, ""}

func (h handlerMap) Group(relativPath string) handlerMap {

	h.PrefixPath = path.Join(h.PrefixPath, relativPath)
	return h
}

func (h handlerMap) HandlerWithRateLimit(method string, relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	(h.Map)[handlerKey{URL: path.Join(h.PrefixPath, relativPath), Method: method}] = handlerProp{
		Handler: handler, NeedRoles: roles, limit: ratelimit, rates: concurrent.NewMap()}
}

func (h handlerMap) GetWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodGet, relativPath, handler, ratelimit, roles...)
}
func (h handlerMap) PostWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodPost, relativPath, handler, ratelimit, roles...)
}
func (h handlerMap) DeleteWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodDelete, relativPath, handler, ratelimit, roles...)
}
func (h handlerMap) PutWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodPut, relativPath, handler, ratelimit, roles...)
}

func (h handlerMap) Handle(method string, relativPath string, handler gin.HandlerFunc, roles ...string) {
	(h.Map)[handlerKey{URL: path.Join(h.PrefixPath, relativPath), Method: method}] = handlerProp{
		Handler: handler, NeedRoles: roles, limit: -1}
}
func (h handlerMap) Get(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodGet, relativPath, handler, roles...)
}
func (h handlerMap) Post(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPost, relativPath, handler, roles...)
}
func (h handlerMap) Delete(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodDelete, relativPath, handler, roles...)
}
func (h handlerMap) Put(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPut, relativPath, handler, roles...)
}
func (h handlerMap) RegisterMapedHandlers(r *gin.Engine) {
	for k, v := range h.Map {
		r.Handle(k.Method, k.URL, v.Handler)
	}
}
