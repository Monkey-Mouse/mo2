// +build normal

package controller

import (
	"testing"
)

func TestController_Log(t *testing.T) {
	req := get(t, "/api"+apiLogs, nil)
	testHTTP(t, tests{name: "Test logs", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
	req.Header.Set("Cookie", "jwtToken=aaaa")
	testHTTP(t, tests{name: "Test bad cookie", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
}
