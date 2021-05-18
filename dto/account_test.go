package dto

import (
	"reflect"
	"testing"

	"github.com/Monkey-Mouse/mo2/server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestContains(t *testing.T) {
	type args struct {
		slice []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test contain", args: args{slice: []string{"xxxx", "xxx"}, item: "xxx"}, want: true},
		{name: "test not contain", args: args{slice: []string{"xxxx", "xxx"}, item: "xxxxx"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapAccount2InfoBrief(t *testing.T) {
	type args struct {
		a model.Account
	}
	id := primitive.NewObjectID()
	tests := []struct {
		name  string
		args  args
		wantB UserInfoBrief
	}{
		{name: "test parse", args: args{model.Account{ID: id, UserName: "user", Settings: map[string]string{"a": "b"}}}, wantB: UserInfoBrief{ID: id, Name: "user", Settings: map[string]string{"a": "b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := MapAccount2InfoBrief(tt.args.a); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("MapAccount2InfoBrief() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
