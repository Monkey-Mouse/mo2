package controller

import (
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

func upsertGroup(t *testing.T, id primitive.ObjectID, ownerID primitive.ObjectID) *http.Request {
	return put(t, "/api/group", nil, model.Group{
		ID:            id,
		OwnerID:       ownerID,
		AccessManager: model.AccessManager{RoleMap: map[string][]primitive.ObjectID{"admin": {ownerID}}},
	})
}
func TestController_UpsertGroup(t *testing.T) {
	groupID := primitive.NewObjectID()

	defer database.DeleteGroupByID(groupID)

	id := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	database.UpsertAccount(&model.Account{
		ID:       id,
		UserName: id.Hex(),
		Email:    id.Hex(),
		Roles:    []string{}})
	defer database.DeleteAccount(id)
	database.UpsertGroup(model.Group{
		ID:            groupID,
		OwnerID:       id,
		AccessManager: model.AccessManager{},
	})

	req1 := upsertGroup(t, groupID, id)
	req2 := upsertGroup(t, groupID, id)
	req3 := upsertGroup(t, groupID, id)
	req4 := upsertGroup(t, groupID, id)
	addCookie(req2)
	addCookieWithID(req3, id3)
	addCookieWithID(req4, id)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test update another group", req: req2, wantCode: 400},
		tests{name: "Test json no bind", req: req2, wantCode: 400},
		tests{name: "Test no user", req: req3, wantCode: 400},
		tests{name: "Test success", req: req4, wantCode: 201, wantStr: "admin"},
	)

}
