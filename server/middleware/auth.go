package middleware

import (
	"mo2/server/controller/badresponse"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modern-go/concurrent"
	"github.com/willf/bloom"
)

var duration int = 10
var unblockEvery int = 3600
var fromCTX FromCTX
var blockFilter = bloom.NewWithEstimates(10000, 0.01)
var userInfoKey string

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
// 请不要手动注册这个中间件，你应该用这个package中的RegisterMapedHandlers方法
func AuthMiddleware(c *gin.Context) {
	// Block illegal ips
	if blockFilter.TestString(c.ClientIP()) {
		c.AbortWithStatusJSON(http.StatusForbidden, badresponse.SetResponseReason("IP Blocked!检测到该ip地址存在潜在的ddos行为"))
		return
	}

	key := handlerKey{c.FullPath(), c.Request.Method}
	prop, ok := handlers[key]
	// not registered for this middleware
	if !ok {
		c.Next()
		return
	}
	// rate limit logic
	if !checkRateLimit(prop, c.ClientIP()) {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, badresponse.SetResponseReason("Too frequent!"))
		return
	}
	// role auth logic
	if prop.NeedRoles == nil || len(prop.NeedRoles) == 0 {
		c.Next()
		return
	}
	uinfo, jwterr := fromCTX(c)
	c.Set(userInfoKey, uinfo)
	if jwterr != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, badresponse.SetResponseReason("Unauthorized!"))
		return
	}
	for _, v := range prop.NeedRoles {
		if !uinfo.IsInRole(v) {
			c.AbortWithStatusJSON(http.StatusForbidden, badresponse.SetResponseReason("Need role: "+v))
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
	Roles      []string
	Limit      int
}

var handlers = make(map[handlerKey]handlerProp, 0)

// H handlermap, like gin router
var H = handlerMap{handlers, "", []string{}, -1}

func (h handlerMap) Group(relativPath string, roles ...string) handlerMap {
	h.PrefixPath = path.Join(h.PrefixPath, relativPath)
	h.Roles = roles
	return h
}
func (h handlerMap) GroupWithLimit(relativPath string, ratelimit int, roles ...string) handlerMap {
	h.PrefixPath = path.Join(h.PrefixPath, relativPath)
	h.Roles = roles
	h.Limit = ratelimit
	return h
}

func (h handlerMap) HandlerWithRateLimit(method string, relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	(h.Map)[handlerKey{URL: path.Join(h.PrefixPath, relativPath), Method: method}] = handlerProp{
		Handler: handler, NeedRoles: append(roles, h.Roles...), limit: ratelimit, rates: concurrent.NewMap()}
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
	h.HandlerWithRateLimit(method, relativPath, handler, h.Limit, roles...)
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

// RegisterMapedHandlers 必须要使用的方法，只有用了它，路由和中间件才会真正被注册
// 使用这个方法请不要手动注册中间件
func (h handlerMap) RegisterMapedHandlers(r *gin.Engine, getUserFromCTX FromCTX, userKey string) {
	fromCTX = getUserFromCTX
	userInfoKey = userKey
	r.Use(AuthMiddleware)
	for k, v := range h.Map {
		r.Handle(k.Method, k.URL, v.Handler)
	}
}
