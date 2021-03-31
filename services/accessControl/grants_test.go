package accessControl

import (
	"github.com/Monkey-Mouse/go-abac/abac"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/dto"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBlog(t *testing.T) {
	isDraft := true
	var blog model.Blog
	// visitor not agree with authorID, but in the group with rights to access
	//managerID:=InsertManager4Test()
	//manager,mErr:=database.FindManager(managerID)
	//if mErr.IsError(){
	//	t.Error(mErr)
	//}
	visitorID := primitive.NewObjectID()
	roleMap := make(map[string][]primitive.ObjectID)
	roleMap["admin"] = []primitive.ObjectID{visitorID}
	manager := model.AccessManager{
		ID:         primitive.NewObjectID(),
		EntityInfo: model.InitEntity(),
		RoleMap:    roleMap,
	}
	if mErr := database.UpsertManager(manager); mErr.IsError() {
		t.Error(mErr)
	}
	defer database.DeleteManagerByID(manager.ID)
	blog.AuthorID = manager.ID
	if mErr := database.UpsertBlog(&blog, isDraft); mErr.IsError() {
		t.Error(mErr)
	}
	defer database.DeleteBlogs(isDraft, blog.ID)
	group := model.Group{
		ID:              primitive.NewObjectID(),
		OwnerID:         visitorID,
		AccessManagerID: manager.ID,
	}
	if mErr := database.UpsertGroup(group); mErr.IsError() {
		t.Error(mErr)
	}
	defer database.DeleteGroupByID(group.ID)
	if pass, err := Ctrl.CanOr(abac.IQueryInfo{
		Subject:  "account",
		Action:   abac.ActionUpdate,
		Resource: "blog",
		Context: abac.DefaultContext{"allowOwn": AllowOwn{ID: blog.ID, Filter: model.Filter{IsDraft: isDraft}, UserInfo: dto.LoginUserInfo{ID: visitorID}},
			"accessFilter": AccessFilter{
				VisitorID: visitorID,
				ManagerID: manager.ID,
				RoleList:  nil,
			}},
	}); err != nil {
		t.Error(err)
	} else {
		if !pass {
			t.Error("not pass")
		}
	}
}
