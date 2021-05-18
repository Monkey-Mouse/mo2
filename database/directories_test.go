package database

import (
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestChooseCol(t *testing.T) {
	tests := []struct {
		name       string
		collection string
		wantCol    *mongo.Collection
	}{
		{"category", CategoryCol, catCol},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCol, gotMErr := ChooseCol(tt.collection)
			if !reflect.DeepEqual(gotCol, tt.wantCol) {
				t.Errorf("ChooseCol() gotCol = %v, want %v", gotCol, tt.wantCol)
			}
			if gotMErr.IsError() {
				t.Errorf("ChooseCol() gotMErr = %v", gotMErr)
			}
		})
	}
}

func TestFindDirectoryInfo(t *testing.T) {
	id := InsertCategories4Test(1)[0]
	name0 := "test0"
	UpsertCategory(&model.Directory{ID: id, Name: name0, ParentID: primitive.NewObjectID()})
	defer DeleteCategoryCompletely(id)
	type args struct {
		collection string
		ids        []primitive.ObjectID
	}
	tests := []struct {
		name     string
		args     args
		wantInfo string
	}{
		{"test", args{
			collection: CategoryCol,
			ids:        []primitive.ObjectID{id},
		}, name0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, gotMErr := FindDirectoryInfo(tt.args.collection, tt.args.ids...)

			if gotMErr.IsError() {
				t.Errorf("ChooseCol() gotMErr = %v", gotMErr)
			}
			if !gotInfo[0].ParentID.IsZero() {
				t.Errorf("ChooseCol() project parentID = %v", gotInfo[0].ParentID)
			}
			if gotInfo[0].Name != tt.wantInfo {
				t.Errorf("ChooseCol() name = %v,want %v", gotInfo[0].Name, tt.wantInfo)
			}
		})
	}
}
