package controller

import (
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

func updateGroup(t *testing.T, id primitive.ObjectID, ownerID primitive.ObjectID) *http.Request {
	return put(t, "/api/group", nil, model.Group{
		ID:            id,
		OwnerID:       ownerID,
		AccessManager: model.AccessManager{RoleMap: map[string][]primitive.ObjectID{"admin": {ownerID}}},
	})
}

func insertGroup(t *testing.T, id primitive.ObjectID, ownerID primitive.ObjectID) *http.Request {
	return post(t, "/api/group", nil, model.Group{
		ID:            id,
		OwnerID:       ownerID,
		AccessManager: model.AccessManager{RoleMap: map[string][]primitive.ObjectID{"admin": {ownerID}}},
	})
}
func deleteGroup(t *testing.T, id primitive.ObjectID) *http.Request {
	return delete(t, "/api/group/"+id.Hex(), nil, nil)
}
func TestController_UpdateGroup(t *testing.T) {
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
	defer database.DeleteGroupByOwnerID(id)
	defer database.DeleteGroupByOwnerID(id3)
	database.UpsertGroup(model.Group{
		ID:            groupID,
		OwnerID:       id,
		AccessManager: model.AccessManager{},
	})

	req1 := updateGroup(t, groupID, id)
	req2 := updateGroup(t, groupID, id)
	req3 := updateGroup(t, groupID, id)
	req4 := updateGroup(t, groupID, id)
	addCookie(req2)
	addCookieWithID(req3, id3)
	addCookieWithID(req4, id)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test update another group", req: req2, wantCode: 400},
		tests{name: "Test no user", req: req3, wantCode: 400},
		tests{name: "Test success", req: req4, wantCode: 201, wantStr: "admin"},
	)

}

func TestController_InsertGroup(t *testing.T) {
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

	req1 := insertGroup(t, groupID, id)
	req2 := insertGroup(t, groupID, id3)
	req3 := insertGroup(t, primitive.ObjectID{}, primitive.ObjectID{})
	req4 := insertGroup(t, primitive.ObjectID{}, primitive.ObjectID{})
	req5 := insertGroup(t, groupID, id)
	addCookieWithID(req2, id)
	addCookieWithID(req3, id3)
	addCookieWithID(req4, id3)
	addCookieWithID(req5, id)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test insert with other's id ", req: req2, wantCode: 201, wantStr: id.Hex()},
		tests{name: "Test json no bind", req: req2, wantCode: 400},
		tests{name: "Test no account id", req: req3, wantCode: 201, wantStr: id3.Hex()},
		tests{name: "Test no group id", req: req4, wantCode: 201, wantStr: id3.Hex()},
		tests{name: "Test success", req: req5, wantCode: 201, wantStr: id.Hex()},
	)
	defer database.DeleteGroupByID(groupID)
	defer database.DeleteGroupByOwnerID(id)
	defer database.DeleteGroupByOwnerID(id3)

}

func TestController_DeleteGroup(t *testing.T) {
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

	req1 := deleteGroup(t, groupID)
	req2 := deleteGroup(t, groupID)
	req3 := deleteGroup(t, primitive.NewObjectID())
	req4 := deleteGroup(t, primitive.ObjectID{})
	req5 := deleteGroup(t, groupID)
	addCookieWithID(req2, id3)
	addCookieWithID(req3, id3)
	addCookieWithID(req4, id3)
	addCookieWithID(req5, id)
	testHTTP(t,
		tests{name: "Test auth", req: req1, wantCode: 403},
		tests{name: "Test delete other's group ", req: req2, wantCode: 400, wantStr: badresponse.NoAccessReason},
		tests{name: "Test group not exist", req: req3, wantCode: 400, wantStr: badresponse.NoAccessReason},
		tests{name: "Test no group id", req: req4, wantCode: 400, wantStr: badresponse.BadRequestReason},
		tests{name: "Test success", req: req5, wantCode: 202},
	)
	defer database.DeleteGroupByID(groupID)
	defer database.DeleteGroupByOwnerID(id)
	defer database.DeleteGroupByOwnerID(id3)

}
