package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mo2/mo2utils/mo2errors"
	"mo2/server/model"
	"reflect"
	"testing"
)

func TestChooseCol(t *testing.T) {
	tests := []struct {
		name       string
		collection string
		wantCol    *mongo.Collection
		//wantMErr mo2errors.Mo2Errors
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

	type args struct {
		collection string
		ids        []primitive.ObjectID
	}
	tests := []struct {
		name     string
		args     args
		wantInfo []model.Directory
		wantMErr mo2errors.Mo2Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, gotMErr := FindDirectoryInfo(tt.args.collection, tt.args.ids...)
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("FindDirectoryInfo() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
			if !reflect.DeepEqual(gotMErr, tt.wantMErr) {
				t.Errorf("FindDirectoryInfo() gotMErr = %v, want %v", gotMErr, tt.wantMErr)
			}
		})
	}
}
