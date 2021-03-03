package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	dto "mo2/dto"
	"mo2/mo2utils"
	"mo2/server/middleware"
	"mo2/server/model"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var r *gin.Engine
var c *Controller

func TestMain(m *testing.M) {
	// Write code here to run before tests
	middleware.SetupRateLimiter(5, 5)
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
func get(t *testing.T, uri string, params map[string]string) (req *http.Request) {
	return send("GET", t, uri, params, nil)
}
func post(t *testing.T, uri string, params map[string]string, body interface{}) (req *http.Request) {
	return send("POST", t, uri, params, body)
}
func put(t *testing.T, uri string, params map[string]string, body interface{}) (req *http.Request) {
	return send("PUT", t, uri, params, body)
}
func delete(t *testing.T, uri string, params map[string]string, body interface{}) (req *http.Request) {
	return send("DELETE", t, uri, params, body)
}
func send(mthd string, t *testing.T, uri string, params map[string]string, body interface{}) (req *http.Request) {
	uri = uri + "?"
	for k, v := range params {
		uri = uri + k + "=" + v + "&"
	}
	v, _ := json.Marshal(body)
	req, err := http.NewRequest(mthd, strings.Trim(uri, "&"), bytes.NewBuffer(v))
	if err != nil {
		t.Fatal(err)
	}
	return
}

func addCookie(req *http.Request) {
	addCookieWithID(req, primitive.NewObjectID())
}
func addCookieWithID(req *http.Request, id primitive.ObjectID) {
	addCookieWithIDAndEmail(req, id, "")
}
func addCookieWithIDAndEmail(req *http.Request, id primitive.ObjectID, email string) {
	req.Header.Set("Cookie",
		"jwtToken="+mo2utils.GenerateJwtCode(dto.LoginUserInfo{Email: email, ID: id, Roles: []string{model.OrdinaryUser}}))
}

type tests struct {
	name        string
	req         *http.Request
	wantCode    int
	wantStr     string
	wantHeaders []string
}

func testHTTP(t *testing.T, testSlice ...tests) {
	for _, test := range testSlice {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, test.req)
			if resp.Code == test.wantCode {
				if p, err := ioutil.ReadAll(resp.Body); err != nil {
					t.Errorf("response err")
				} else if !strings.Contains(string(p), test.wantStr) {
					t.Errorf("Want contain str: %v, actual: %v", test.wantStr, string(p))
				} else {
					for _, v := range test.wantHeaders {
						_, ok := resp.HeaderMap[v]
						if !ok {
							t.Errorf("Want header: %v", v)
						}
					}
				}
			} else {
				t.Errorf("Want code: %v, actual code: %v", test.wantCode, resp.Code)
			}
		})
	}
}
