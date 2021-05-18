package controller

import (
	"net/http"
	"os"
	"testing"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
func addAccount(t *testing.T, name string, email string) *http.Request {
	return post(t, "/api/accounts", nil, model.AddAccount{UserName: name, Email: email, Password: "xxxxaaaa"})
}
func deleteAccount(t *testing.T, pass string, email string) *http.Request {
	return delete(t, "/api/accounts", nil, model.DeleteAccount{Email: email, Password: pass})
}
func login(t *testing.T, pass string, emailOrName string) *http.Request {
	return post(t, "/api/accounts/login", nil, model.LoginAccount{UserNameOrEmail: emailOrName, Password: pass})
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
		tests{name: "Test json no bind", req: req1, wantCode: 400})
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
		tests{name: "Test json no bind", req: req2, wantCode: 400},
		tests{name: "Test no user", req: req3, wantCode: 404},
		tests{name: "Test success", req: req4, wantCode: 200, wantStr: "hi"},
	)
	database.DeleteAccount(id)
}

func TestController_AddAccount(t *testing.T) {
	id := primitive.NewObjectID()
	id1 := primitive.NewObjectID()
	database.UpsertAccount(&model.Account{
		ID:       id,
		UserName: id.Hex(),
		Email:    id.Hex(),
		Roles:    []string{},
		Infos:    map[string]string{model.IsActive: model.True},
	})
	database.UpsertAccount(&model.Account{
		ID:       id1,
		UserName: id1.Hex(),
		Email:    id1.Hex(),
		Roles:    []string{},
		Infos:    map[string]string{model.IsActive: model.False},
	})
	req1 := addAccount(t, "", "lll")
	req2 := addAccount(t, "", "lll")
	req3 := addAccount(t, id.Hex(), "lll")
	req4 := addAccount(t, "lll", id.Hex())
	req5 := addAccount(t, id1.Hex(), "lll")
	req6 := addAccount(t, "lll", id1.Hex())
	addCookie(req2)
	addCookie(req3)
	addCookie(req4)
	addCookie(req5)
	addCookie(req6)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test empty account name", req: req2, wantCode: 422},
		tests{name: "Test json no bind", req: req2, wantCode: 400},
		tests{name: "Test name dup", req: req3, wantCode: 422, wantStr: "Name"},
		tests{name: "Test email dup", req: req4, wantCode: 422, wantStr: "Email"},
		tests{name: "Test name dup inactive", req: req5, wantCode: 200},
		tests{name: "Test email dup inactive", req: req6, wantCode: 200},
	)
	database.DeleteAccount(id)
	database.DeleteAccount(id1)
}

func TestController_DeleteAccount(t *testing.T) {
	id := primitive.NewObjectID()
	pass := "aaaaaaaa"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)
	database.UpsertAccount(&model.Account{
		ID:        id,
		UserName:  id.Hex(),
		Email:     id.Hex(),
		Roles:     []string{},
		Infos:     map[string]string{model.IsActive: model.True},
		HashedPwd: string(hash),
	})
	req1 := deleteAccount(t, id.Hex(), id.Hex())
	req2 := deleteAccount(t, id.Hex(), id.Hex())
	req3 := deleteAccount(t, id.Hex(), id.Hex())
	req4 := deleteAccount(t, pass, id.Hex())
	addCookie(req2)
	addCookieWithIDAndEmail(req3, id, id.Hex())
	addCookieWithIDAndEmail(req4, id, id.Hex())
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test delete others", req: req2, wantCode: 422},
		tests{name: "Test json not bind", req: req2, wantCode: 400},
		tests{name: "Test wrong pass", req: req3, wantCode: 403},
		tests{name: "Test delete", req: req4, wantCode: 204},
	)
	_, err := database.DeleteAccount(id)
	if !err.IsError() {
		t.Errorf("Account didn't really delete!")
	}
}

func TestController_LoginAccount(t *testing.T) {
	id := primitive.NewObjectID()
	id1 := primitive.NewObjectID()
	pass := "aaaaaaaa"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)
	database.UpsertAccount(&model.Account{
		ID:        id,
		UserName:  id.Hex(),
		Email:     id.Hex(),
		Roles:     []string{},
		Infos:     map[string]string{model.IsActive: model.True},
		HashedPwd: string(hash),
	})
	database.UpsertAccount(&model.Account{
		ID:        id1,
		UserName:  id1.Hex(),
		Email:     id1.Hex(),
		Roles:     []string{},
		Infos:     map[string]string{model.IsActive: model.False},
		HashedPwd: string(hash),
	})
	req1 := login(t, "xxxx", id1.Hex())
	req2 := login(t, pass, id1.Hex())
	req3 := login(t, "xxx", id.Hex())
	req4 := login(t, "", "")
	req5 := login(t, pass, id.Hex())
	addCookie(req2)
	addCookie(req3)
	addCookie(req4)
	addCookie(req5)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test not activate", req: req2, wantCode: 401},
		tests{name: "Test wrong pass", req: req3, wantCode: 401},
		tests{name: "Test empty data", req: req4, wantCode: 422},
		tests{name: "Test login", req: req5, wantCode: 200, wantHeaders: []string{"Set-Cookie"}},
	)

	database.DeleteAccount(id)
	database.DeleteAccount(id1)
}
