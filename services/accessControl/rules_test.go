package accessControl

import (
	"context"
	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

// InsertManager4Test 插入n个category，并返回它们的id列表
func InsertManager4Test() (id primitive.ObjectID) {
	roleMap := make(map[string][]primitive.ObjectID)
	roleMap["admin"] = []primitive.ObjectID{primitive.NewObjectID()}

	res, err := AccessManagerCol.InsertOne(context.TODO(), model.AccessManager{
		ID:         primitive.NewObjectID(),
		EntityInfo: model.InitEntity(),
		RoleMap:    roleMap,
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
		RoleList  []string
	}
	managerID := InsertManager4Test()
	adminID := primitive.NewObjectID()
	AccessManagerCol.UpdateOne(context.TODO(), bson.M{"_id": managerID}, bson.M{"$push": bson.M{"role_map.admin": adminID}})
	var manager model.AccessManager
	if err := AccessManagerCol.FindOne(context.TODO(), bson.M{"_id": managerID}).Decode(&manager); err != nil {
		t.Error(err)
	}
	log.Println(manager)
	defer func() {
		if _, err := AccessManagerCol.DeleteOne(context.TODO(), bson.M{"_id": managerID}); err != nil {
			t.Error(err)
		}
	}()

	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{

		{"foo", fields{VisitorID: adminID, ManagerID: managerID, RoleList: []string{"admin"}}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &accessFilter{
				VisitorID: tt.fields.VisitorID,
				ManagerID: tt.fields.ManagerID,
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
