package middleware

import (
	"errors"
	"mo2/server/controller/badresponse"
	"net/http"
	"path"
	"strings"
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
	if prop.needRoles == nil || len(prop.needRoles) == 0 {
		c.Next()
		return
	}
	uinfo, jwterr := fromCTX(c)
	c.Set(userInfoKey, uinfo)
	if jwterr != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, badresponse.SetResponseReason("Unauthorized!"))
		return
	}
	if err := checkRoles(uinfo, prop.needRoles); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, badresponse.SetResponseError(err))
		return
	}
	c.Next()
}

func checkRoles(uinfo RoleHolder, rolePolicies [][]string) error {
	for _, v := range rolePolicies {
		failedCheck := true
		for _, u := range v {
			if uinfo.IsInRole(u) {
				failedCheck = false
				break
			}
		}
		if failedCheck {
			return errors.New("Need role: " + strings.Join(v, " or "))
		}
	}
	return nil
}

type handlerProp struct {
	handler   gin.HandlerFunc
	needRoles [][]string
	rates     *concurrent.Map
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

var handlers = make(map[handlerKey]handlerProp, 0)

// H handlermap, like gin router
var H = handlerMap{handlers, "", make([][]string, 0), -1}

// Group 类似gin router的Group方法，注意group里设置的多个role是以or逻辑连接的
// 而group下属的其它api设置的role条件会与group里的role条件进行与运算
func (h handlerMap) Group(relativePath string, roles ...string) handlerMap {
	h.prefixPath = path.Join(h.prefixPath, relativePath)
	h.roles = append(h.roles, roles)
	return h
}

// GroupWithLimit 字面意思，增加了ratelimit功能
func (h handlerMap) GroupWithLimit(relativePath string, ratelimit int, roles ...string) handlerMap {
	h = h.Group(relativePath, roles...)
	h.limit = ratelimit
	return h
}

// HandlerWithRateLimit 用于非常规http方法的处理
func (h handlerMap) HandlerWithRateLimit(method string, relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.roles = append(h.roles, roles)
	(h.innerMap)[handlerKey{url: path.Join(h.prefixPath, relativPath), method: method}] = handlerProp{
		handler: handler, needRoles: h.roles, limit: ratelimit, rates: concurrent.NewMap()}
}

// GetWithRateLimit 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) GetWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodGet, relativPath, handler, ratelimit, roles...)
}

// PostWithRateLimit 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) PostWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodPost, relativPath, handler, ratelimit, roles...)
}

// DeleteWithRateLimit 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) DeleteWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodDelete, relativPath, handler, ratelimit, roles...)
}

// PutWithRateLimit 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) PutWithRateLimit(relativPath string, handler gin.HandlerFunc, ratelimit int, roles ...string) {
	h.HandlerWithRateLimit(http.MethodPut, relativPath, handler, ratelimit, roles...)
}

// Handler 用于非常规http方法的处理
func (h handlerMap) Handle(method string, relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.HandlerWithRateLimit(method, relativPath, handler, h.limit, roles...)
}

// Get 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Get(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodGet, relativPath, handler, roles...)
}

// Post 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Post(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPost, relativPath, handler, roles...)
}

// Delete 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Delete(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodDelete, relativPath, handler, roles...)
}

// Put 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Put(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(http.MethodPut, relativPath, handler, roles...)
}

// RegisterMapedHandlers 必须要使用的方法，只有用了它，路由和中间件才会真正被注册
// 使用这个方法请不要手动注册中间件
func (h handlerMap) RegisterMapedHandlers(r *gin.Engine, getUserFromCTX FromCTX, userKey string) {
	fromCTX = getUserFromCTX
	userInfoKey = userKey
	r.Use(AuthMiddleware)
	for k, v := range h.innerMap {
		r.Handle(k.method, k.url, v.handler)
	}
}
