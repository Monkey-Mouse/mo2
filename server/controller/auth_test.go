// +build rate_limit

package controller

import (
	"mo2/mo2utils"
	"mo2/server/middleware"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestHandlers(c *Controller) {
	api := middleware.H.Group("/apitest")
	{
		api.GetWithRL(apiLogs, c.Log, 10)
	}
}

func Test_AuthMiddleware(t *testing.T) {
	authR := gin.New()
	req := get(t, "/apitest"+apiLogs, nil)
	c = NewController()
	setupTestHandlers(c)
	middleware.H.RegisterMapedHandlers(authR, func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
		return
	}, mo2utils.UserInfoKey)
	ch := make(chan bool, 0)
	t.Run("Test rate limit block", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			go func() {
				resp := httptest.NewRecorder()
				authR.ServeHTTP(resp, req)
				ch <- resp.Code == 200
			}()
		}
		sucNum := 0
		for i := 0; i < 100; i++ {
			success := <-ch
			if success {
				sucNum++
			}
		}
		if sucNum < 10 || sucNum > 20 {
			t.Errorf("auth middleware should ban 80-90 times, actual baned: %v", 100-sucNum)
		}
	})
	time.Sleep(5 * time.Second)
	t.Run("Test rate limit unblock", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			go func() {
				resp := httptest.NewRecorder()
				authR.ServeHTTP(resp, req)
				ch <- resp.Code == 200
			}()
		}
		sucNum := 0
		for i := 0; i < 100; i++ {
			success := <-ch
			if success {
				sucNum++
			}
		}
		if sucNum < 10 || sucNum > 20 {
			t.Errorf("auth middleware should ban 80-90 times, actual baned: %v", 100-sucNum)
		}
	})
}
