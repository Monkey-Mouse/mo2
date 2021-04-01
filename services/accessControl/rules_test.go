package accessControl

import (
	"context"
	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

// InsertGroup4Test 插入n个category，并返回它们的id列表
func InsertGroup4Test() (id primitive.ObjectID) {
	roleMap := make(map[string][]primitive.ObjectID)
	roleMap["admin"] = []primitive.ObjectID{primitive.NewObjectID()}
	manager := model.AccessManager{
		EntityInfo: model.InitEntity(),
		RoleMap:    roleMap,
	}
	res, err := database.GroupCol.InsertOne(context.TODO(), model.Group{
		ID:            primitive.NewObjectID(),
		OwnerID:       primitive.ObjectID{},
		AccessManager: manager,
	})
	if err != nil {
		log.Fatal(err)
	}
	id = res.InsertedID.(primitive.ObjectID)
	return
}
func Test_accessFilter_JudgeRule(t *testing.T) {
	type fields struct {
		VisitorID primitive.ObjectID
		ManagerID primitive.ObjectID
		RoleList  [][]string
	}
	groupID := InsertGroup4Test()
	adminID := primitive.NewObjectID()
	database.GroupCol.UpdateOne(context.TODO(), bson.M{"_id": groupID}, bson.M{"$push": bson.M{"access_manager.role_map.admin": adminID}})
	var group model.Group
	if err := database.GroupCol.FindOne(context.TODO(), bson.M{"_id": groupID}).Decode(&group); err != nil {
		t.Error(err)
	}
	log.Println(group)
	defer func() {
		if _, err := database.GroupCol.DeleteOne(context.TODO(), bson.M{"_id": groupID}); err != nil {
			t.Error(err)
		}
	}()

	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{

		{"foo", fields{VisitorID: adminID, ManagerID: groupID, RoleList: [][]string{{"admin"}}}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccessFilter{
				VisitorID: tt.fields.VisitorID,
				GroupID:   tt.fields.ManagerID,
				RoleList:  tt.fields.RoleList,
			}
			got, err := a.JudgeRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("JudgeRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JudgeRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}
