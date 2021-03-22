package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/modern-go/concurrent"
	"github.com/willf/bloom"
)

const redisHashSet = "ratelimiter"

var duration int = 10
var unblockEvery int = 3600
var fromCTX FromCTX
var blockFilter = bloom.NewWithEstimates(10000, 0.01)
var userInfoKey string
var handlers = make(map[handlerKey]handlerProp, 0)
var rdb *redis.Client
var dicChan chan *concurrent.Map = make(chan *concurrent.Map, 0)
var lock = sync.Mutex{}

// H handlermap, like gin router
var H = handlerMap{handlers, "", make([][]string, 0), -1}

func resetVar() {
	blockFilter = bloom.NewWithEstimates(10000, 0.01)
	handlers = make(map[handlerKey]handlerProp, 0)
	dicChan = make(chan *concurrent.Map, 0)
	lock = sync.Mutex{}
	H = handlerMap{handlers, "", make([][]string, 0), -1}
	rdb = nil
}

// OptionalParams config rate limiter
type OptionalParams struct {
	LimitEvery   int
	Unblockevery int
	UseRedis     bool
}

// setupRateLimiter setup ddos banner
func setupRateLimiter(limitEvery int, unblockevery int, useRedis bool) {
	if useRedis {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			db = 0
		}
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASS"), // no password set
			DB:       db,                      // use default DB
		})
		rdb.Del(redisHashSet).Err()
	}
	duration = limitEvery
	unblockEvery = unblockevery
	blockFilter.ClearAll()
}
func redisCleaner() {
	for {
		time.Sleep(time.Second * time.Duration(duration))
		if rdb == nil {
			return
		}
		rdb.Del(redisHashSet).Err()
	}

}
func cleaner() {
	cancelChan := make(chan bool, 0)
	for {
		go func() {
			newDic := concurrent.NewMap()
			for {
				select {
				case <-cancelChan:
					return
				case dicChan <- newDic:
				}
			}
		}()
		time.Sleep(time.Second * time.Duration(duration))
		cancelChan <- true
	}
}
func resetBlocker() {
	for {
		time.Sleep(time.Second * time.Duration(unblockEvery))
		blockFilter.ClearAll()
	}
}

func redisCheckRL(prop string, ip string, limit int) bool {
	// rate limit logic
	if limit < 0 {
		return true
	}
	key := prop + ip
	re, err := rdb.HIncrBy(redisHashSet, key, 1).Result()
	if err != nil {
		log.Fatal(err)
	}
	if re > int64(limit) {
		blockFilter.AddString(ip)
		return false
	}
	return true
}

func checkRL(prop string, ip string, limit int) bool {
	// rate limit logic
	if limit < 0 {
		return true
	}
	key := prop + ip
	dic := <-dicChan
	lock.Lock()
	v, ext := dic.Load(key)
	if !ext {
		dic.Store(key, 1)
	} else {
		dic.Store(key, (v.(int))+1)
		if v.(int)+1 > limit {
			blockFilter.AddString(ip)
			lock.Unlock()
			return false
		}
	}
	lock.Unlock()
	return true
}

func checkBlock(ip string) *mo2errors.Mo2Errors {
	if blockFilter.TestString(ip) {
		return mo2errors.New(http.StatusForbidden, "IP Blocked!检测到该ip地址存在潜在的ddos行为")
	}
	return nil
}

func checkBlockAndRL(prop string, ip string, limit int) *mo2errors.Mo2Errors {
	err := checkBlock(ip)
	if err != nil {
		return err
	}
	var ok bool
	if rdb != nil {
		ok = redisCheckRL(prop, ip, limit)
	} else {
		ok = checkRL(prop, ip, limit)
	}
	if !ok {
		return mo2errors.New(http.StatusTooManyRequests, "Too frequent!")
	}
	return nil
}

// AuthMiddleware also have rate limit function
// 请不要手动注册这个中间件，你应该用这个package中的RegisterMapedHandlers方法
func AuthMiddleware(c *gin.Context) {
	key := handlerKey{c.FullPath(), c.Request.Method}
	prop, ok := handlers[key]
	// not registered for this middleware
	if !ok {
		c.Next()
		return
	}
	// rate limit logic
	err := checkBlockAndRL(key.method+key.url, c.ClientIP(), prop.limit)
	if err != nil {
		c.AbortWithStatusJSON(err.ErrorCode, err.ErrorTip)
	}
	uinfo, jwterr := fromCTX(c)
	c.Set(userInfoKey, uinfo)
	// role auth logic
	if prop.needRoles == nil || len(prop.needRoles) == 0 {
		c.Next()
		return
	}
	if jwterr != nil {
		badresponse.SetErrResponse(c, http.StatusForbidden,
			"Unauthorized!")
		return
	}
	if err := checkRoles(uinfo, prop.needRoles); err != nil {
		badresponse.SetErrResponse(c, http.StatusForbidden,
			err.Error())
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

// Group 类似gin router的Group方法，注意group里设置的多个role是以or逻辑连接的
// 而group下属的其它api设置的role条件会与group里的role条件进行与运算
func (h handlerMap) Group(
	relativePath string,
	roles ...string) handlerMap {
	h.prefixPath = path.Join(h.prefixPath, relativePath)
	if roles != nil && len(roles) != 0 {
		h.roles = append(h.roles, roles)
	}
	return h
}

// GroupWithRL 字面意思，增加了ratelimit功能
func (h handlerMap) GroupWithRL(
	relativePath string,
	ratelimit int,
	roles ...string) handlerMap {
	h = h.Group(relativePath, roles...)
	h.limit = ratelimit
	return h
}

// HandlerWithRL 用于非常规http方法的处理
func (h handlerMap) HandlerWithRL(
	method string,
	relativPath string,
	handler gin.HandlerFunc,
	ratelimit int,
	roles ...string) {
	if roles != nil && len(roles) != 0 {
		h.roles = append(h.roles, roles)
	}
	key := handlerKey{
		url:    path.Join(h.prefixPath, relativPath),
		method: method,
	}
	(h.innerMap)[key] = handlerProp{
		handler:   handler,
		needRoles: h.roles,
		limit:     ratelimit,
	}
}

// GetWithRL 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) GetWithRL(
	relativPath string,
	handler gin.HandlerFunc,
	ratelimit int,
	roles ...string) {
	h.HandlerWithRL(
		http.MethodGet,
		relativPath,
		handler,
		ratelimit,
		roles...,
	)
}

// PostWithRL 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) PostWithRL(
	relativPath string,
	handler gin.HandlerFunc,
	ratelimit int,
	roles ...string) {
	h.HandlerWithRL(
		http.MethodPost,
		relativPath,
		handler,
		ratelimit,
		roles...,
	)
}

// DeleteWithRL 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) DeleteWithRL(
	relativPath string,
	handler gin.HandlerFunc,
	ratelimit int,
	roles ...string) {
	h.HandlerWithRL(
		http.MethodDelete,
		relativPath,
		handler,
		ratelimit,
		roles...,
	)
}

// PutWithRL 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) PutWithRL(
	relativPath string,
	handler gin.HandlerFunc,
	ratelimit int,
	roles ...string) {
	h.HandlerWithRL(
		http.MethodPut,
		relativPath,
		handler,
		ratelimit,
		roles...,
	)
}

// Handler 用于非常规http方法的处理
func (h handlerMap) Handle(
	method string,
	relativPath string,
	handler gin.HandlerFunc,
	roles ...string) {
	h.HandlerWithRL(
		method,
		relativPath,
		handler,
		h.limit,
		roles...,
	)
}

// Get 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Get(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(
		http.MethodGet,
		relativPath,
		handler,
		roles...,
	)
}

// Post 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Post(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(
		http.MethodPost,
		relativPath,
		handler,
		roles...,
	)
}

// Delete 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Delete(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(
		http.MethodDelete,
		relativPath,
		handler,
		roles...,
	)
}

// Put 方法接收的role条件以或逻辑连接，连接之后的逻辑与父亲的group里的逻辑以与逻辑链接
func (h handlerMap) Put(relativPath string, handler gin.HandlerFunc, roles ...string) {
	h.Handle(
		http.MethodPut,
		relativPath,
		handler,
		roles...,
	)
}

// RegisterMapedHandlers 必须要使用的方法，只有用了它，路由和中间件才会真正被注册
// 使用这个方法请不要手动注册中间件
func (h handlerMap) RegisterMapedHandlers(
	r *gin.Engine,
	getUserFromCTX FromCTX,
	userKey string,
	optional *OptionalParams) {
	if optional != nil {
		setupRateLimiter(optional.LimitEvery, optional.Unblockevery, optional.UseRedis)
	}
	fromCTX = getUserFromCTX
	userInfoKey = userKey
	r.Use(AuthMiddleware)

	for k, v := range h.innerMap {
		r.Handle(k.method, k.url, v.handler)
	}
	if rdb == nil {
		go cleaner()
	} else {
		go redisCleaner()
	}

	go resetBlocker()
}
