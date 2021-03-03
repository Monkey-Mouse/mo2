package controller

import (
	"mo2/server/model"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestController_Log(t *testing.T) {
	req := get(t, "/api"+apiLogs, nil)
	testHTTP(t, tests{name: "Test logs", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
	req.Header.Set("Cookie", "jwtToken=aaaa")
	testHTTP(t, tests{name: "Test bad cookie", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
}

func TestController_AddAccountRole(t *testing.T) {
	req := post(t, "/api/accounts/role", nil, model.AddAccountRole{ID: primitive.NewObjectID(), Roles: []string{"a"}, SuperKey: "xx"})
	req1 := req.Clone(req.Context())
	addCookie(req1)
	testHTTP(t,
		tests{name: "Test auth", req: req, wantCode: 403},
		tests{name: "Test add account role", req: req1, wantCode: 401})
}
