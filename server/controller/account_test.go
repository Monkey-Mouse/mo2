package controller

import (
	"mo2/database"
	"mo2/server/model"
	"net/http"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestController_Log(t *testing.T) {
	req := get(t, "/api"+apiLogs, nil)
	testHTTP(t, tests{name: "Test logs", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
	req.Header.Set("Cookie", "jwtToken=aaaa")
	testHTTP(t, tests{name: "Test bad cookie", req: req, wantCode: 200, wantStr: "", wantHeaders: []string{"Set-Cookie"}})
}
func addRoleReq(t *testing.T, id primitive.ObjectID) *http.Request {
	return post(t, "/api/accounts/role", nil, model.AddAccountRole{ID: id, Roles: []string{"a"}, SuperKey: "xx"})
}
func TestController_AddAccountRole(t *testing.T) {
	id := primitive.NewObjectID()
	req := addRoleReq(t, id)
	req1 := addRoleReq(t, id)
	req2 := addRoleReq(t, id)
	req3 := addRoleReq(t, id)
	addCookie(req1)
	addCookie(req2)
	addCookie(req3)
	testHTTP(t,
		tests{name: "Test auth", req: req, wantCode: 403},
		tests{name: "Test add account role no right", req: req1, wantCode: 401})
	os.Setenv("MO2_SUPER_KEY", "xx")
	testHTTP(t,
		tests{name: "Test add account role no user", req: req2, wantCode: 404})
	database.UpsertAccount(&model.Account{ID: id, UserName: id.Hex(), Email: id.Hex(), Roles: []string{}})
	testHTTP(t,
		tests{name: "Test add account role", req: req3, wantCode: 200})
	database.DeleteAccount(id)
}
