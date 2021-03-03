package controller

import (
	"mo2/database"
	"mo2/dto"
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
func updateAccount(t *testing.T, id primitive.ObjectID) *http.Request {
	return put(t, "/api/accounts", nil, dto.UserInfoBrief{ID: id, Name: "xx", Settings: map[string]string{"hi": "hello"}})
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
		tests{name: "Test add account role no right", req: req1, wantCode: 401},
		tests{name: "Test json no bind", req: req1, wantCode: 401})
	os.Setenv("MO2_SUPER_KEY", "xx")
	testHTTP(t,
		tests{name: "Test add account role no user", req: req2, wantCode: 404})
	database.UpsertAccount(&model.Account{ID: id, UserName: id.Hex(), Email: id.Hex(), Roles: []string{}})
	testHTTP(t,
		tests{name: "Test add account role", req: req3, wantCode: 200})
	database.DeleteAccount(id)
}

func TestController_UpdateAccount(t *testing.T) {
	id := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	database.UpsertAccount(&model.Account{
		ID:       id,
		UserName: id.Hex(),
		Email:    id.Hex(),
		Roles:    []string{}})
	req1 := updateAccount(t, id)
	req2 := updateAccount(t, id)
	req3 := updateAccount(t, id3)
	req4 := updateAccount(t, id)
	addCookie(req2)
	addCookieWithID(req3, id3)
	addCookieWithID(req4, id)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test update another account", req: req2, wantCode: 403},
		tests{name: "Test json no bind", req: req2, wantCode: 401},
		tests{name: "Test no user", req: req3, wantCode: 404},
		tests{name: "Test success", req: req4, wantCode: 200, wantStr: "hi"})
	database.DeleteAccount(id)
}
