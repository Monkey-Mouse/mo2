package controller

import (
	"mo2/mo2utils"
	"mo2/server/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var c *Controller

func TestMain(m *testing.M) {
	// Write code here to run before tests
	r = gin.Default()
	c = NewController()
	SetupHandlers(c)
	middleware.H.RegisterMapedHandlers(r, func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
		str, err := ctx.Cookie("jwtToken")
		if err != nil {
			return
		}
		userInfo, err = mo2utils.ParseJwt(str)
		return
	}, mo2utils.UserInfoKey)
	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}
func get(t *testing.T, uri string, params map[string]string) (req *http.Request, resp *httptest.ResponseRecorder) {
	resp = httptest.NewRecorder()
	uri = uri + "?"
	for k, v := range params {
		uri = uri + k + "=" + v + "&"
	}

	req, err := http.NewRequest("GET", strings.Trim(uri, "&"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return
}
